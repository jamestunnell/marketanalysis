import van from "vanjs-core"

import { ButtonIconTooltip } from '../elements/buttons'
import Checkbox from '../elements/checkbox'
import { MakeConstraint } from '../constraint.js'
import { IconCheck, IconError } from '../elements/icons'
import { IntRange, FloatRange } from '../elements/number'
import { Table, TableRow } from '../elements/table'

const {tbody} = van.tags

class TargetParam {
    constructor({address, constraint, step, defaultValue, makeValueInput}) {
        this.address = address
        this.constraint = constraint
        this.step = step
        this.makeValueInput = makeValueInput

        const min = constraint.getMin()
        let minVal = defaultValue
        if (min) {
            minVal = min.inclusive ? min.value : (min.value + step)
        } else if (max) {
            minVal = max.value - step
        }

        const max = constraint.getMax()
        let maxVal = defaultValue
        if (max) {
            maxVal = max.inclusive ? max.value : (max.value - step)
        } else if (max) {
            maxVal = min.value + step
        }

        this.selected = van.state(false)
        this.min = van.state(minVal)
        this.max = van.state(maxVal)
        this.err = van.derive(() => {
            const minErr = this.constraint.validate(this.min.val)
            if (minErr) {
                return new Error(`min value ${this.min.val} is invalid: ${minErr.message}`)
            }

            const maxErr = this.constraint.validate(this.max.val)
            if (maxErr) {
                return new Error(`max value ${this.max.val} is invalid: ${maxErr.message}`)
            }

            if (this.min.val >= this.max.val) {
                return new Error(`min value ${this.min.val} is not < max value ${this.max.val}`)
            }

            return null
        })
    }

    renderRow() {
        const checkbox = Checkbox({checked: this.selected})
        const min = this.makeValueInput({
            constraint: this.constraint,
            value: this.min,
        })
        const max = this.makeValueInput({
            constraint: this.constraint,
            value: this.max,
        })
        const status = ButtonIconTooltip({
            icon: () => this.err.val ? IconError() : IconCheck(),
            text: () => this.err.val ? `Invalid min/max value: ${this.err.val.message}` : "Min and max values are valid"
        })

        return TableRow([this.address, checkbox, min, max, status])
    }
}

const MakeTargetParam = ({address, paramInfo}) => {
    const constraint = MakeConstraint(paramInfo.constraint)

    // don't support enum param yet
    if (!constraint.isRange()) {
        return null
    }

    switch (paramInfo.valueType) {
        case 'int':
            console.log("making int target param", {address, paramInfo})

            return new TargetParam({
                address, constraint,
                step:1,
                defaultValue: paramInfo.defaultValue,
                makeValueInput: IntRange,
            })
        case 'float64':
            console.log("making float64 target param", {address, paramInfo})

            return new TargetParam({
                address, constraint,
                step:0.01,
                defaultValue: paramInfo.defaultValue,
                makeValueInput: FloatRange,
            })
    }

    return null
}

const TargetParamsTable = (targetParams) => {
    return Table({
        columnNames: ["Parameter", "Optimize", "Min", "Max", ""],
        tableBody: tbody({class:"table-auto"}, targetParams.map(t => t.renderRow())),
    })
}

export { MakeTargetParam, TargetParamsTable }