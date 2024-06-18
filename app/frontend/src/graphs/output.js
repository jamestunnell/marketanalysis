import van from "vanjs-core"

import Checkbox from '../elements/checkbox.js'
import { Table, TableRow } from '../elements/table.js'

const { tbody } = van.tags

const MEASURE_FIRST = "first"
const MEASURE_LAST = "last"
const MEASURE_MEAN = "mean"
const MEASURE_MIN = "min"
const MEASURE_MAX = "max"
const MEASURE_STDDEV = "stddev"

class Output {
    constructor({name, measurements}) {
        this.name = name
        this.measurements = measurements
    }

    isMeasurementsEmpty() {
        return this.measurements.val.length == 0
    }

    renderTableRow() {
        const firstChecked = van.state(this.measurements.val.indexOf(MEASURE_FIRST) >= 0)
        const lastChecked = van.state(this.measurements.val.indexOf(MEASURE_LAST) >= 0)
        const minChecked = van.state(this.measurements.val.indexOf(MEASURE_MIN) >= 0)
        const maxChecked = van.state(this.measurements.val.indexOf(MEASURE_MAX) >= 0)
        const meanChecked = van.state(this.measurements.val.indexOf(MEASURE_MEAN) >= 0)
        const stddevChecked = van.state(this.measurements.val.indexOf(MEASURE_STDDEV) >= 0)
    
        van.derive(() => {
            const newMeasurements = []
    
            if (firstChecked.val) {
                newMeasurements.push(MEASURE_FIRST)
            }
    
            if (lastChecked.val) {
                newMeasurements.push(MEASURE_LAST)
            }

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
    
            this.measurements.val = newMeasurements
        })
    
        return TableRow([
            this.name,
            Checkbox({checked: firstChecked}),
            Checkbox({checked: lastChecked}),
            Checkbox({checked: minChecked}),
            Checkbox({checked: maxChecked}),
            Checkbox({checked: meanChecked}),
            Checkbox({checked: stddevChecked}),
        ])
    }

    makeConfig() {
        return {name: this.name, measurements: this.measurements.val}
    }
}

const OutputsTable = ({outputs, hidden}) => {
    return Table({
        hidden,
        columnNames: ["Name", "First", "Last", "Min", "Max", "Mean", "Std. Dev."],
        tableBody: tbody({class:"table-auto"}, outputs.map(output => output.renderTableRow())),
    })
}

function MakeOutputs({infos, configs}) {
    return infos.map(info => {
        const cfg = configs.find(cfg => cfg.name === info.name)
        
        return new Output({
            name: info.name,
            measurements: van.state(cfg ? cfg.measurements : []),
        })
    })
}

export {OutputsTable, MakeOutputs}