import van from "vanjs-core"

import { Table, TableRow } from "../table"
import { INPUT_CLASS } from "../input"

const { input, tbody} = van.tags

const MEASURE_MEAN = "mean"
const MEASURE_MIN = "min"
const MEASURE_MAX = "max"
const MEASURE_STDDEV = "stddev"

const Checkbox = (checked) => {
    return input({
        class: INPUT_CLASS,
        type: "checkbox",
        onchange: e => checked.val = e.target.checked,
        checked: checked.val,
    })
}

const OutputRow = ({name, measurements}) => {
    const minChecked = van.state(measurements.val.indexOf(MEASURE_MIN) >= 0)
    const maxChecked = van.state(measurements.val.indexOf(MEASURE_MAX) >= 0)
    const meanChecked = van.state(measurements.val.indexOf(MEASURE_MEAN) >= 0)
    const stddevChecked = van.state(measurements.val.indexOf(MEASURE_STDDEV) >= 0)

    van.derive(() => {
        const newMeasurements = []

        if (minChecked.val) {
            newMeasurements.push(MEASURE_MIN)
        }

        if (maxChecked.val) {
            newMeasurements.push(MEASURE_MAX)
        }
 
        if (meanChecked.val) {
            newMeasurements.push(MEASURE_MEAN)
        }
        
        if (stddevChecked.val) {
            newMeasurements.push(MEASURE_STDDEV)
        }

        measurements.val = newMeasurements
    })

    return TableRow([
        name,
        Checkbox(minChecked),
        Checkbox(maxChecked),
        Checkbox(meanChecked),
        Checkbox(stddevChecked),
    ])    
}

const OutputsTable = ({outputs, measurements}) => {
    const names = outputs.map(output => output.name).sort()

    const rows = names.map(name => {
        return OutputRow({name, measurements: measurements[name]})
    })

    return Table({
        columnNames: ["Name", "Min", "Max", "Mean", "Std. Dev."],
        tableBody: tbody({class:"table-auto"}, rows),
    })
}

export { OutputsTable }