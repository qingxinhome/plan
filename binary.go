package main

var (
	//+
	//date + interval
	gBinDateIntervalAdd binDateInterAddOp

	// *
	gBinFloat32Multi binFloat32MultiOp

	// =
	gBinInt32Equal binInt32EqualOp

	// >
	gBinInt32Great   binInt32GreatOp
	gBinFloat32Great binFloat32GreatOp

	gBinDateIntervalSingleOpWrapper   binarySingleOpWrapper[Date, Interval, Date]
	gBinInt32BoolSingleOpWrapper      binarySingleOpWrapper[int32, int32, bool]
	gBinFloat32Float32SingleOpWrapper binarySingleOpWrapper[float32, float32, float32]
	gBinFloat32BoolSingleOpWrapper    binarySingleOpWrapper[float32, float32, bool]
)

type binaryOp[T any, S any, R any] interface {
	operation(left *T, right *S, result *R)
}

type binaryFunc[T any, S any, R any] interface {
	fun(left *T, right *S, result *R)
}

type binaryWrapper[T any, S any, R any] interface {
	operation(left *T, right *S, result *R, mask *Bitmap, idx int,
		op binaryOp[T, S, R],
		fun binaryFunc[T, S, R])

	addsNulls() bool
}

type binarySingleOpWrapper[T any, S any, R any] struct {
}

func (b binarySingleOpWrapper[T, S, R]) operation(left *T, right *S, result *R, mask *Bitmap, idx int,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R]) {
	op.operation(left, right, result)
}

func (b binarySingleOpWrapper[T, S, R]) addsNulls() bool {
	return false
}

// +
type binDateInterAddOp struct {
}

func (op binDateInterAddOp) operation(left *Date, right *Interval, result *Date) {
	if right._unit == "year" {
		result._year = left._year + right._year
		result._month = left._month
		result._day = left._day
	} else {
		panic("usp")
	}
}

// *
type binFloat32MultiOp struct{}

func (op binFloat32MultiOp) operation(left, right *float32, result *float32) {
	*result = *left * *right
}

// = int32
type binInt32EqualOp struct {
}

func (op binInt32EqualOp) operation(left, right *int32, result *bool) {
	*result = *left == *right
}

// > int32
type binInt32GreatOp struct {
}

func (op binInt32GreatOp) operation(left, right *int32, result *bool) {
	*result = *left > *right
}

// > float32
type binFloat32GreatOp struct {
}

func (op binFloat32GreatOp) operation(left, right *float32, result *bool) {
	*result = *left > *right
}

func binaryExecSwitch[T any, S any, R any](
	left, right, result *Vector,
	count int,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R],
	wrapper binaryWrapper[T, S, R],
) {
	if left.phyFormat().isConst() && right.phyFormat().isConst() {
		binaryExecConst[T, S, R](left, right, result, count, op, fun, wrapper)
	} else if left.phyFormat().isFlat() && right.phyFormat().isConst() {
		binaryExecFlat[T, S, R](left, right, result, count, op, fun, wrapper, false, true)
	} else if left.phyFormat().isConst() && right.phyFormat().isFlat() {
		binaryExecFlat[T, S, R](left, right, result, count, op, fun, wrapper, true, false)
	} else if left.phyFormat().isFlat() && right.phyFormat().isFlat() {
		binaryExecFlat[T, S, R](left, right, result, count, op, fun, wrapper, false, false)
	} else {
		binaryExecGeneric[T, S, R](left, right, result, count, op, fun, wrapper)
	}
}
func binaryExecConst[T any, S any, R any](
	left, right, result *Vector,
	count int,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R],
	wrapper binaryWrapper[T, S, R],
) {
	result.setPhyFormat(PF_CONST)
	if isNullInPhyFormatConst(left) ||
		isNullInPhyFormatConst(right) {
		setNullInPhyFormatConst(result, true)
		return
	}
	lSlice := getSliceInPhyFormatConst[T](left)
	rSlice := getSliceInPhyFormatConst[S](right)
	resSlice := getSliceInPhyFormatConst[R](result)

	wrapper.operation(&lSlice[0], &rSlice[0], &resSlice[0], getMaskInPhyFormatConst(result), 0, op, fun)
}

func binaryExecFlat[T any, S any, R any](
	left, right, result *Vector,
	count int,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R],
	wrapper binaryWrapper[T, S, R],
	lconst, rconst bool,
) {
	lSlice := getSliceInPhyFormatFlat[T](left)
	rSlice := getSliceInPhyFormatFlat[S](right)
	if lconst && isNullInPhyFormatConst(left) ||
		rconst && isNullInPhyFormatConst(right) {
		result.setPhyFormat(PF_CONST)
		setNullInPhyFormatConst(result, true)
		return
	}

	result.setPhyFormat(PF_FLAT)
	resSlice := getSliceInPhyFormatFlat[R](result)
	resMask := getMaskInPhyFormatFlat(result)
	if lconst {
		if wrapper.addsNulls() {
			resMask.copyFrom(getMaskInPhyFormatFlat(right), count)
		} else {
			setMaskInPhyFormatFlat(result, getMaskInPhyFormatFlat(right))
		}
	} else if rconst {
		if wrapper.addsNulls() {
			resMask.copyFrom(getMaskInPhyFormatFlat(left), count)
		} else {
			setMaskInPhyFormatFlat(result, getMaskInPhyFormatFlat(left))
		}
	} else {
		if wrapper.addsNulls() {
			resMask.copyFrom(getMaskInPhyFormatFlat(left), count)
			if resMask.AllValid() {
				resMask.copyFrom(getMaskInPhyFormatFlat(right), count)
			} else {
				resMask.combine(getMaskInPhyFormatFlat(right), count)
			}
		} else {
			setMaskInPhyFormatFlat(result, getMaskInPhyFormatFlat(left))
			resMask.combine(getMaskInPhyFormatFlat(right), count)
		}
	}
	binaryExecFlatLoop[T, S, R](
		lSlice,
		rSlice,
		resSlice,
		count,
		resMask,
		op,
		fun,
		wrapper,
		lconst,
		rconst,
	)
}

func binaryExecFlatLoop[T any, S any, R any](
	ldata []T, rdata []S,
	resData []R,
	count int,
	mask *Bitmap,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R],
	wrapper binaryWrapper[T, S, R],
	lconst, rconst bool,
) {
	if !mask.AllValid() {
		baseIdx := 0
		eCnt := entryCount(count)
		for i := 0; i < eCnt; i++ {
			ent := mask.getEntry(uint64(i))
			next := min(baseIdx+8, count)
			if AllValidInEntry(ent) {
				for ; baseIdx < next; baseIdx++ {
					lidx := baseIdx
					ridx := baseIdx
					if lconst {
						lidx = 0
					}
					if rconst {
						ridx = 0
					}
					wrapper.operation(&ldata[lidx], &rdata[ridx], &resData[baseIdx], mask, baseIdx, op, fun)
				}
			} else if NoneValidInEntry(ent) {
				baseIdx = next
				continue
			} else {
				start := baseIdx
				for ; baseIdx < next; baseIdx++ {
					if rowIsValidInEntry(ent, uint64(baseIdx-start)) {
						lidx := baseIdx
						ridx := baseIdx
						if lconst {
							lidx = 0
						}
						if rconst {
							ridx = 0
						}
						wrapper.operation(&ldata[lidx], &rdata[ridx], &resData[baseIdx], mask, baseIdx, op, fun)
					}
				}
			}
		}
	} else {
		for i := 0; i < count; i++ {
			lidx := i
			ridx := i
			if lconst {
				lidx = 0
			}
			if rconst {
				ridx = 0
			}
			wrapper.operation(&ldata[lidx], &rdata[ridx], &resData[i], mask, i, op, fun)
		}
	}
}

func binaryExecGeneric[T any, S any, R any](
	left, right, result *Vector,
	count int,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R],
	wrapper binaryWrapper[T, S, R],
) {
	var ldata, rdata *UnifiedFormat
	left.toUnifiedFormat(count, ldata)
	right.toUnifiedFormat(count, rdata)

	lSlice := getSliceInPhyFormatUnifiedFormat[T](ldata)
	rSlice := getSliceInPhyFormatUnifiedFormat[S](rdata)
	result.setPhyFormat(PF_FLAT)
	resSlice := getSliceInPhyFormatFlat[R](result)
	binaryExecGenericLoop[T, S, R](
		lSlice,
		rSlice,
		resSlice,
		ldata._sel,
		rdata._sel,
		count,
		ldata._mask,
		rdata._mask,
		result._mask,
		op,
		fun,
		wrapper,
	)
}

func binaryExecGenericLoop[T any, S any, R any](
	ldata []T, rdata []S,
	resData []R,
	lsel *SelectVector,
	rsel *SelectVector,
	count int,
	lmask *Bitmap,
	rmask *Bitmap,
	resMask *Bitmap,
	op binaryOp[T, S, R],
	fun binaryFunc[T, S, R],
	wrapper binaryWrapper[T, S, R],
) {
	if !lmask.AllValid() || !rmask.AllValid() {
		for i := 0; i < count; i++ {
			lidx := lsel.getIndex(i)
			ridx := rsel.getIndex(i)
			if lmask.rowIsValid(uint64(lidx)) && rmask.rowIsValid(uint64(ridx)) {
				wrapper.operation(&ldata[lidx], &rdata[ridx], &resData[i], resMask, i, op, fun)
			} else {
				resMask.setInvalid(uint64(i))
			}
		}
	} else {
		for i := 0; i < count; i++ {
			lidx := lsel.getIndex(i)
			ridx := rsel.getIndex(i)
			wrapper.operation(&ldata[lidx], &rdata[ridx], &resData[i], resMask, i, op, fun)
		}
	}
}