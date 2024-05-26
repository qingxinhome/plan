package main

import (
	"unsafe"
)

func StateSize[T any, STATE State[T]]() int {
	var val STATE
	size := unsafe.Sizeof(val)
	return int(size)
}

func UnaryAggregate[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	inputTyp LType,
	retTyp LType,
	nullHandling FuncNullHandling,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) *AggrFunc {
	var size aggrStateSize
	var init aggrInit
	var update aggrUpdate
	var combine aggrCombine
	var finalize aggrFinalize
	var simpleUpdate aggrSimpleUpdate
	size = func() int {
		var val State[ResultT]
		return int(unsafe.Sizeof(val))
	}
	init = func(pointer unsafe.Pointer) {
		aop.Init((*State[ResultT])(pointer), sop)
	}
	update = func(inputs []*Vector, data *AggrInputData, inputCount int, states *Vector, count int) {
		assertFunc(inputCount == 1)
		UnaryScatter[ResultT, STATE, InputT, OP](inputs[0], states, data, count, aop, sop, addOp, top)
	}
	combine = func(source *Vector, target *Vector, data *AggrInputData, count int) {
		Combine[ResultT, STATE, InputT, OP](source, target, data, count, aop, sop, addOp, top)
	}
	finalize = func(states *Vector, data *AggrInputData, result *Vector, count int, offset int) {
		Finalize[ResultT, STATE, InputT, OP](states, data, result, count, offset, aop, sop, addOp, top)
	}
	simpleUpdate = func(inputs []*Vector, data *AggrInputData, inputCount int, state unsafe.Pointer, count int) {
		assertFunc(inputCount == 1)
		UnaryUpdate[ResultT, STATE, InputT, OP](inputs[0], data, state, count, aop, sop, addOp, top)
	}
	return &AggrFunc{
		_args:         []LType{inputTyp},
		_retType:      retTyp,
		_stateSize:    size,
		_init:         init,
		_update:       update,
		_combine:      combine,
		_finalize:     finalize,
		_nullHandling: nullHandling,
		_simpleUpdate: simpleUpdate,
	}
}

func GetSumAggr(pTyp PhyType) *AggrFunc {
	switch pTyp {
	case INT32:
		fun := UnaryAggregate[Hugeint, State[Hugeint], int32, SumOp[Hugeint, int32]](
			integer(),
			hugeint(),
			DEFAULT_NULL_HANDLING,
			SumOp[Hugeint, int32]{},
			&SumStateOp[Hugeint]{},
			&HugeintAdd{},
			&Hugeint{},
		)
		return fun
	default:
		panic("usp")
	}
}

type TypeOp[T any] interface {
	Add(*T, *T)
	Mul(*T, *T)
}

type StateType int

const (
	STATE_SUM StateType = iota
)

type State[T any] struct {
	_typ   StateType
	_isset bool
	_value T
}

func (state *State[T]) Init() {
	switch state._typ {
	case STATE_SUM:
		state._isset = false
	default:
		panic("usp")
	}

}

func (state *State[T]) Combine(other *State[T], add TypeOp[T]) {
	switch state._typ {
	case STATE_SUM:
		state._isset = other._isset || state._isset
		add.Add(&state._value, &other._value)
	default:
		panic("usp")
	}

}

func (state *State[T]) SetIsset(b bool) {
	state._isset = b
}
func (state *State[T]) SetValue(val T) {
	state._value = val
}

func (state *State[T]) GetIsset() bool {
	return state._isset
}

func (state *State[T]) GetValue() T {
	return state._value
}

type StateOp[T any] interface {
	Init(*State[T])
	Combine(*State[T], *State[T], *AggrInputData, TypeOp[T])
	AddValues(*State[T], int)
}

type AddOp[ResultT any, InputT any] interface {
	AddNumber(*State[ResultT], *InputT, TypeOp[ResultT])
	AddConstant(*State[ResultT], *InputT, int, TypeOp[ResultT])
}

type AggrOp[ResultT any, InputT any] interface {
	Init(*State[ResultT], StateOp[ResultT])
	Combine(*State[ResultT], *State[ResultT], *AggrInputData, StateOp[ResultT], TypeOp[ResultT])
	Operation(*State[ResultT], *InputT, *AggrUnaryInput, StateOp[ResultT], AddOp[ResultT, InputT], TypeOp[ResultT])
	ConstantOperation(*State[ResultT], *InputT, *AggrUnaryInput, int, StateOp[ResultT], AddOp[ResultT, InputT], TypeOp[ResultT])
	Finalize(*State[ResultT], *ResultT, *AggrFinalizeData)
	IgnoreNull() bool
}

type SumStateOp[T any] struct {
}

func (SumStateOp[T]) Init(s *State[T]) {
	s.Init()
}
func (SumStateOp[T]) Combine(src *State[T], target *State[T], _ *AggrInputData, top TypeOp[T]) {
	src.Combine(target, top)
}
func (SumStateOp[T]) AddValues(s *State[T], _ int) {
	s.SetIsset(true)
}

type HugeintAdd struct {
}

func (*HugeintAdd) addValue(result *Hugeint, value uint64, positive int) {
	result._lower += value
	overflow := 0
	if result._lower < value {
		overflow = 1
	}
	if overflow^positive == 0 {
		result._upper += -1 + 2*int64(positive)
	}
}

func (hadd *HugeintAdd) AddNumber(state *State[Hugeint], input *int32, top TypeOp[Hugeint]) {
	pos := 0
	if *input >= 0 {
		pos = 1
	}
	hadd.addValue(&state._value, uint64(*input), pos)
}

func (*HugeintAdd) AddConstant(*State[Hugeint], *int32, int, TypeOp[Hugeint]) {
	//TODO:
}

type SumOp[ResultT any, InputT any] struct {
}

func (s SumOp[ResultT, InputT]) Init(s2 *State[ResultT], sop StateOp[ResultT]) {
	var val ResultT
	s2.SetValue(val)
	sop.Init(s2)
}

func (s SumOp[ResultT, InputT]) Combine(src *State[ResultT], target *State[ResultT], data *AggrInputData,
	sop StateOp[ResultT], top TypeOp[ResultT]) {
	sop.Combine(src, target, data, top)
}

func (s SumOp[ResultT, InputT]) Operation(s3 *State[ResultT], input *InputT, data *AggrUnaryInput,
	sop StateOp[ResultT], aop AddOp[ResultT, InputT], top TypeOp[ResultT]) {
	sop.AddValues(s3, 1)
	aop.AddNumber(s3, input, top)
}

func (s SumOp[ResultT, InputT]) ConstantOperation(s3 *State[ResultT], input *InputT, data *AggrUnaryInput, count int,
	sop StateOp[ResultT], aop AddOp[ResultT, InputT], top TypeOp[ResultT]) {
	sop.AddValues(s3, count)
	aop.AddConstant(s3, input, count, top)
}

func (s SumOp[ResultT, InputT]) Finalize(s3 *State[ResultT], target *ResultT, data *AggrFinalizeData) {
	if !s3.GetIsset() {
		data.ReturnNull()
	} else {
		*target = s3.GetValue()
	}
}

func (s SumOp[ResultT, InputT]) IgnoreNull() bool {
	return true
}

func UnaryScatter[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	input *Vector,
	states *Vector,
	data *AggrInputData,
	count int,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	if input.phyFormat().isConst() &&
		states.phyFormat().isConst() {
		if aop.IgnoreNull() && isNullInPhyFormatConst(input) {
			return
		}
		inputSlice := getSliceInPhyFormatConst[InputT](input)
		statesPtrSlice := getSliceInPhyFormatConst[unsafe.Pointer](states)
		inputData := NewAggrUnaryInput(data, getMaskInPhyFormatConst(input))
		aop.ConstantOperation((*State[ResultT])(statesPtrSlice[0]), &inputSlice[0], inputData, count, sop, addOp, top)
	} else if input.phyFormat().isFlat() && states.phyFormat().isFlat() {
		inputSlice := getSliceInPhyFormatFlat[InputT](input)
		statesPtrSlice := getSliceInPhyFormatFlat[unsafe.Pointer](states)
		UnaryFlatLoop[ResultT, STATE, InputT, OP](
			inputSlice,
			data,
			statesPtrSlice,
			getMaskInPhyFormatFlat(input),
			count,
			aop,
			sop,
			addOp,
			top,
		)
	} else {
		var idata, sdata UnifiedFormat
		input.toUnifiedFormat(count, &idata)
		states.toUnifiedFormat(count, &sdata)
		UnaryScatterLoop[ResultT, STATE, InputT](
			getSliceInPhyFormatUnifiedFormat[InputT](&idata),
			data,
			getSliceInPhyFormatUnifiedFormat[unsafe.Pointer](&sdata),
			idata._sel,
			sdata._sel,
			idata._mask,
			count,
			aop,
			sop,
			addOp,
			top,
		)
	}
}

func UnaryFlatLoop[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	inputSlice []InputT,
	data *AggrInputData,
	statesPtrSlice []unsafe.Pointer,
	mask *Bitmap,
	count int,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	if aop.IgnoreNull() && !mask.AllValid() {
		input := NewAggrUnaryInput(data, mask)
		baseIdx := &input._inputIdx
		*baseIdx = 0
		eCnt := entryCount(count)
		for eIdx := 0; eIdx < eCnt; eIdx++ {
			e := mask.getEntry(uint64(eIdx))
			next := min(*baseIdx+8, count)
			if AllValidInEntry(e) {
				for ; *baseIdx < next; *baseIdx++ {
					aop.Operation((*State[ResultT])(statesPtrSlice[*baseIdx]), &inputSlice[*baseIdx], input, sop, addOp, top)
				}
			} else if NoneValidInEntry(e) {
				*baseIdx = next
				continue
			} else {
				start := *baseIdx
				for ; *baseIdx < next; *baseIdx++ {
					if rowIsValidInEntry(e, uint64(*baseIdx-start)) {
						aop.Operation((*State[ResultT])(statesPtrSlice[*baseIdx]), &inputSlice[*baseIdx], input, sop, addOp, top)
					}
				}
			}
		}
	} else {
		input := NewAggrUnaryInput(data, mask)
		i := &input._inputIdx
		for *i = 0; *i < count; *i++ {
			aop.Operation((*State[ResultT])(statesPtrSlice[*i]), &inputSlice[*i], input, sop, addOp, top)
		}
	}
}

func UnaryScatterLoop[ResultT any, STATE State[ResultT], InputT any](
	inputSlice []InputT,
	data *AggrInputData,
	statesPtrSlice []unsafe.Pointer,
	isel *SelectVector,
	ssel *SelectVector,
	mask *Bitmap,
	count int,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	if aop.IgnoreNull() && !mask.AllValid() {
		input := NewAggrUnaryInput(data, mask)
		for i := 0; i < count; i++ {
			input._inputIdx = isel.getIndex(i)
			sidx := ssel.getIndex(i)
			if mask.rowIsValid(uint64(input._inputIdx)) {
				aop.Operation((*State[ResultT])(statesPtrSlice[sidx]), &inputSlice[input._inputIdx], input, sop, addOp, top)
			}
		}
	} else {
		input := NewAggrUnaryInput(data, mask)
		for i := 0; i < count; i++ {
			input._inputIdx = isel.getIndex(i)
			sidx := ssel.getIndex(i)
			aop.Operation((*State[ResultT])(statesPtrSlice[sidx]), &inputSlice[input._inputIdx], input, sop, addOp, top)
		}
	}
}

func Combine[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	source *Vector,
	target *Vector,
	data *AggrInputData,
	count int,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	assertFunc(source.typ().isPointer())
	assertFunc(target.typ().isPointer())
	sourcePtrSlice := getSliceInPhyFormatFlat[unsafe.Pointer](source)
	targetPtrSlice := getSliceInPhyFormatFlat[unsafe.Pointer](target)
	for i := 0; i < count; i++ {
		aop.Combine((*State[ResultT])(sourcePtrSlice[i]),
			(*State[ResultT])(targetPtrSlice[i]),
			data,
			sop,
			top,
		)
	}
}

func Finalize[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	states *Vector,
	data *AggrInputData,
	result *Vector,
	count int,
	offset int,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	if states.phyFormat().isConst() {
		result.setPhyFormat(PF_CONST)
		statePtrSlice := getSliceInPhyFormatFlat[unsafe.Pointer](states)
		resultSlice := getSliceInPhyFormatFlat[ResultT](result)
		final := NewAggrFinalizeData(result, data)
		aop.Finalize((*State[ResultT])(statePtrSlice[0]), &resultSlice[0], final)
	} else {
		assertFunc(states.phyFormat().isFlat())
		result.setPhyFormat(PF_FLAT)
		statePtrSlice := getSliceInPhyFormatFlat[unsafe.Pointer](states)
		resultSlice := getSliceInPhyFormatFlat[ResultT](result)
		final := NewAggrFinalizeData(result, data)
		for i := 0; i < count; i++ {
			final._resultIdx = i + offset
			aop.Finalize((*State[ResultT])(statePtrSlice[i]), &resultSlice[final._resultIdx], final)
		}
	}
}

func UnaryUpdate[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	input *Vector,
	data *AggrInputData,
	statePtr unsafe.Pointer,
	count int,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	switch input.phyFormat() {
	case PF_CONST:
		if aop.IgnoreNull() && isNullInPhyFormatConst(input) {
			return
		}
		inputSlice := getSliceInPhyFormatFlat[InputT](input)
		inputData := NewAggrUnaryInput(data, getMaskInPhyFormatConst(input))
		aop.ConstantOperation((*State[ResultT])(statePtr), &inputSlice[0], inputData, count, sop, addOp, top)
	case PF_FLAT:
		inputSlice := getSliceInPhyFormatFlat[InputT](input)
		UnaryFlatUpdateLoop[ResultT, STATE, InputT, OP](
			inputSlice,
			data,
			statePtr,
			count,
			getMaskInPhyFormatFlat(input),
			aop,
			sop,
			addOp,
			top,
		)
	default:
		var idata UnifiedFormat
		input.toUnifiedFormat(count, &idata)
		UnaryUpdateLoop[ResultT, STATE, InputT, OP](
			getSliceInPhyFormatUnifiedFormat[InputT](&idata),
			data,
			statePtr,
			count,
			idata._mask,
			idata._sel,
			aop,
			sop,
			addOp,
			top,
		)
	}
}

func UnaryFlatUpdateLoop[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	inputSlice []InputT,
	data *AggrInputData,
	statePtr unsafe.Pointer,
	count int,
	mask *Bitmap,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	input := NewAggrUnaryInput(data, mask)
	baseIdx := &input._inputIdx
	*baseIdx = 0
	eCnt := entryCount(count)
	for eIdx := 0; eIdx < eCnt; eIdx++ {
		e := mask.getEntry(uint64(eIdx))
		next := min(*baseIdx+8, count)
		if !aop.IgnoreNull() || AllValidInEntry(e) {
			for ; *baseIdx < next; *baseIdx++ {
				aop.Operation((*State[ResultT])(statePtr), &inputSlice[*baseIdx], input, sop, addOp, top)
			}
		} else if NoneValidInEntry(e) {
			*baseIdx = next
			continue
		} else {
			start := *baseIdx
			for ; *baseIdx < next; *baseIdx++ {
				if rowIsValidInEntry(e, uint64(*baseIdx-start)) {
					aop.Operation((*State[ResultT])(statePtr), &inputSlice[*baseIdx], input, sop, addOp, top)
				}
			}
		}
	}
}

func UnaryUpdateLoop[ResultT any, STATE State[ResultT], InputT any, OP AggrOp[ResultT, InputT]](
	inputSlice []InputT,
	data *AggrInputData,
	statePtr unsafe.Pointer,
	count int,
	mask *Bitmap,
	selVec *SelectVector,
	aop AggrOp[ResultT, InputT],
	sop StateOp[ResultT],
	addOp AddOp[ResultT, InputT],
	top TypeOp[ResultT],
) {
	input := NewAggrUnaryInput(data, mask)
	if aop.IgnoreNull() && !mask.AllValid() {
		for i := 0; i < count; i++ {
			input._inputIdx = selVec.getIndex(i)
			if mask.rowIsValid(uint64(input._inputIdx)) {
				aop.Operation((*State[ResultT])(statePtr), &inputSlice[input._inputIdx], input, sop, addOp, top)
			}
		}
	} else {
		for i := 0; i < count; i++ {
			input._inputIdx = selVec.getIndex(i)
			aop.Operation((*State[ResultT])(statePtr), &inputSlice[input._inputIdx], input, sop, addOp, top)
		}
	}
}
