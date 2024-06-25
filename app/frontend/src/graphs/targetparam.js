import van from "vanjs-core"

import { ButtonIconTooltip } from '../elements/buttons'
import Checkbox from '../elements/checkbox'
import { MakeConstraint } from '../constraint.js'
import { IconCheck, IconError } from '../elements/icons'
import { IntRange, FloatRange } from '../elements/number'
import { Table, TableRow } from '../elements/table'

const {div, tbody} = van.tags

class EnumTargetParam {
    constructor({address, valueType, constraint}) {
        this.address = address
        this.valueType = valueType
        this.constraint = constraint
        this.selected = van.state(false)
        this.selectedValues = constraint.limits.map(value => van.state(true))
    }

    getAddress() {
        return this.address
    }

    getNewConstraint() {
        return {
            type: `OneOf[${this.valueType}]`,
            limits: this.constraint.limits.filter((value, idx) => this.selectedValues[idx].val),
        }
    }

    renderCheckbox() {
        return Checkbox({checked: this.selected})
    }

    renderLimitsArea() {
        const rows = TableRow(this.constraint.limits.map((value,idx) => {
            return [span(value), Checkbox({checked: this.selectedValues[idx]})]
        }))

        return Table({
            hidden: van.derive(() => !this.selected.val),
            columnNames: ["Value", "Selected"],
            tableBody: tbody({class:"table-auto"}, rows),
        })
    }
}

class RangeTargetParam {
    constructor({address, valueType, constraint, step, defaultValue, makeValueInput}) {
        this.address = address
        this.valueType = valueType
        this.constraint = constraint
        this.step = step
        this.makeValueInput = makeValueInput

        const min = constraint.getMin()
        const max = constraint.getMax()

        let minVal = defaultValue
        if (min) {
            minVal = min.inclusive ? min.value : (min.value + step)
        } else if (max) {
            minVal = max.value - step
        }

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

            if (this.min.val > this.max.val) {
                return new Error(`min value ${this.min.val} is not <= max value ${this.max.val}`)
            }

            return null
        })
    }

    getAddress() {
        return this.address
    }

    getNewConstraint() {
        return {
            type: `RangeIncl[${this.valueType}]`,
            limits: [this.min.val, this.max.val],
        }
    }

    renderCheckbox() {
        return Checkbox({checked: this.selected})
    }

    renderLimitsArea() {
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

        return div(
            {
                class: "flex flex-row",
                hidden: van.derive(() => !this.selected.val)
            },
            min,
            max,
            status,
        )
    }
}

const MakeTargetParam = ({address, paramInfo}) => {
    const constraint = MakeConstraint(paramInfo.constraint)
    const valueType = paramInfo.valueType

    if (!constraint.isRange()) {
        console.log("making enum target param", {address, valueType, paramInfo})

        return new EnumTargetParam({address, valueType, constraint})
    }

    switch (valueType) {
        case 'int':
            console.log("making int range target param", {address, valueType, paramInfo})

            return new RangeTargetParam({
                address, constraint, valueType,
                step:1,
                defaultValue: paramInfo.defaultValue,
                makeValueInput: IntRange,
            })
        case 'float64':
            console.log("making float64 range target param", {address, valueType, paramInfo})

            return new RangeTargetParam({
                address, constraint, valueType,
                step:0.01,
                defaultValue: paramInfo.defaultValue,
                makeValueInput: FloatRange,
            })
    }

    return null
}

const TargetParamsTable = (targetParams) => {
    console.log("making target params table", {targetParams})

    const rows = targetParams.map(t => {
        return TableRow([t.getAddress(), t.renderCheckbox(), t.renderLimitsArea()])
    })
    return Table({
        columnNames: ["Parameter", "Optimize", "Limits"],
        tableBody: tbody({class:"table-auto"}, rows),
    })
}

export { MakeTargetParam, TargetParamsTable }