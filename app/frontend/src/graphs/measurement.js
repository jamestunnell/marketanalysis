const MEASURE_FIRST = "first"
const MEASURE_LAST = "last"
const MEASURE_MEAN = "mean"
const MEASURE_MIN = "min"
const MEASURE_MAX = "max"
const MEASURE_STDDEV = "stddev"

function allMeasurements() {
    return [MEASURE_FIRST, MEASURE_LAST, MEASURE_MEAN, MEASURE_MIN, MEASURE_MAX, MEASURE_STDDEV]
}

export {allMeasurements, MEASURE_FIRST, MEASURE_LAST, MEASURE_MEAN, MEASURE_MIN, MEASURE_MAX, MEASURE_STDDEV}