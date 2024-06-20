import van from "vanjs-core"

import { ButtonIconTooltip } from '../elements/buttons.js'
import { MakeConstraint } from '../constraint.js'
import { IconCheck, IconError } from '../elements/icons.js';
import { Table, TableRow } from '../elements/table.js'
import { IntRange, IntEnum, FloatRange, FloatEnum } from '../elements/number.js'

const { tbody } = van.tags

class Param {
    constructor({constraint, defaultValue, name, value, makeValueInput}) {
        this.constraint = constraint
        this.defaultValue = defaultValue
        this.name = name
        this.value = value
        this.makeValueInput = makeValueInput
    }
    
    isValueDefault() {
        return this.value.val === this.defaultValue
    }

    makeConfig() {
        return {name: this.name, value: this.value.val}
    }

    renderTableRow() {
        const err = van.derive(() => this.constraint.validate(this.value.val))
        const valueInput = this.makeValueInput({
            constraint: this.constraint,
            id: this.name,
            value: this.value,
        })
        const valueStatus = ButtonIconTooltip({
            icon: () => err.val ? IconError() : IconCheck(),
            text: van.derive(() => err.val ? `Value is invalid: ${err.val.message}` : "Value is valid"),
        })
        
        return TableRow([this.name, this.constraint.toString(), valueInput, valueStatus ])
    }
}

function ParamsTable({params, hidden}) {
    return Table({
        hidden,  
        columnNames: ["Name", "Constraint", "Value", ""],
        tableBody: tbody({class:"table-auto"}, params.map(param => param.renderTableRow())),
    })
}

function MakeParams({infos, configs}) {
    const paramVals = Object.fromEntries(configs.map(cfg => {
        return [cfg.name, cfg.value]
    }))
    
    return infos.map(info => {
        const name = info.name
        const value = van.state(paramVals[info.name] ?? info.defaultValue)
        const defaultValue = info.defaultValue
        const constraint = MakeConstraint(info.constraint)

        if (info.valueType == 'int') {
            if (constraint.isRange()) {
                return new Param({name, value, defaultValue, constraint, makeValueInput: IntRange})
            } else {
                return new Param({name, value, defaultValue, constraint, makeValueInput: IntEnum}) 
            }
        } else if (info.valueType == 'float64') {
            if (constraint.isRange()) {
                return new Param({name, value, defaultValue, constraint, makeValueInput: FloatRange})
            } else {
                return new Param({name, value, defaultValue, constraint, makeValueInput: FloatEnum}) 
            }
        } else {
            console.log(`unsupported param type ${info.valueType}`)
        }

        return null
    }).filter(p => p)
}

export {ParamsTable, MakeParams}