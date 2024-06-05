import van from "vanjs-core"

import capitalize from './capitalize.js';
import { ModalBackground } from "./modal.js";
import { Button, ButtonCancel, ButtonIconTooltip } from "./buttons.js";
import { ButtonGroup } from "./buttongroup.js";
import { IconCheck, IconError } from "./icons.js";
import { Table, TableRow } from "./table.js";

const { div, input, label, option, p, select, tbody} = van.tags

const inputClass = "block border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
const minSteps = 100.0

function computeStep(min, max) {
    const x = Math.log10((max - min) / minSteps)
    
    return Math.pow(10.0, Math.ceil(x - 1.0))
}

const EnterNumberRow = ({name, currentVal, min, max, step, updateVal, errMsg}) => {
    const labelText = `${capitalize(name)}: (${min}-${max}):`;

    const status = ButtonIconTooltip({
        icon: () => (errMsg.val.length > 0) ? IconError() : IconCheck(),
        tooltipText: van.derive(() => (errMsg.val.length > 0) ? `Value is invalid: ${errMsg.val}` : "Value is valid"),
    });
    const rowItems = [
        label({for: name}, labelText),
        input({
            id: name,
            type: "number",
            class: inputClass,
            value: currentVal,
            min,
            max,
            step,
            onchange: (e) => {
                const err = updateVal(e.target.value)
                if (err) {
                    errMsg.val = err.message
                } else {
                    errMsg.val = ""
                }
            },
        }),
        status,
    ]

    return TableRow(rowItems)
}

const SelectValueRow = ({name, allowedVals, currentVal, updateVal}) => {
    const options = allowedVals.map(allowedVal => {
        let props = {value: allowedVal};
        
        if (allowedVal === currentVal) {
            props.selected = "selected";
        }

        return option(props, allowedVal);
    });
    // const status = ButtonIconTooltip({
    //     icon: () => (errMsg.val.length > 0) ? IconError() : IconCheck(),
    //     tooltipText: van.derive(() => (errMsg.val.length > 0) ? `Value is invalid: ${errMsg.val}` : "Value is valid"),
    // });

    const rowItems = [
        label({for: name}, capitalize(name)),
        select({
            id: name,
            class: inputClass,
            oninput: (e) => {
                const err = updateVal(e.target.value)
                if (err) {
                    errMsg.val = err.message
                } else {
                    errMsg.val = ""
                }
            },
        }, options),
        // status,
    ]

    return TableRow(rowItems)
}

const IntParamRow = ({param, value, errMsg}) => {
    const updateVal = (strVal) => {
        const newVal = parseInt(strVal)
        if (isNaN(newVal)) {
            return new Error(`${strVal} is not an integer`)
        }
    
        value.val = newVal
    
        return validateParamVal(param, newVal)
    }
    const name = param.name
    const step = 1
    const currentVal = value.val

    let min
    let max

    switch (param.constraint.type) {
        case 'oneOf':
            return SelectValueRow({currentVal, name, updateVal, allowedVals: param.constraint.limits})
        case 'less':
            min = Number.MIN_SAFE_INTEGER
            max = param.constraint.limits[0]-1

            return EnterNumberRow({name, currentVal, step, max, updateVal})
        case 'lessEqual':
            min = Number.MIN_SAFE_INTEGER
            max = param.constraint.limits[0]

            return EnterNumberRow({name, currentVal, step, max, updateVal})
        case 'greater':
            min = param.constraint.limits[0]+1
            max = Number.MAX_SAFE_INTEGER

            return EnterNumberRow({name, currentVal, step, min, updateVal})
        case 'greaterEqual':
            min = param.constraint.limits[0]
            max = Number.MAX_SAFE_INTEGER

            return EnterNumberRow({name, currentVal, step, min, updateVal})
        case 'rangeIncl':
            min = param.constraint.limits[0]
            max = param.constraint.limits[1]

            return EnterNumberRow({name, currentVal, step, min, max, updateVal})
        case 'rangeExcl':
            min = param.constraint.limits[0]
            max = param.constraint.limits[1]-1

            return EnterNumberRow({name, currentVal, step, min, max, updateVal})
    }

    console.log(`unsupported int constraint type ${param.constraint.type}`)

    return null
} 

const FloatParamRow = ({param, value, errMsg}) => {
    const updateVal = (strVal) => {
        const newVal = parseFloat(strVal)
        if (isNaN(newVal)) {
            return new Error(`${strVal} is not an float`)
        }

        value.val = newVal

        const err = validateParamVal(param, newVal)
        if (err) {
            errMsg.val = err.message
        }
    }

    const name = param.name
    const currentVal = value.val

    let min
    let max
    let step
    
    switch (param.constraint.type) {
        case 'oneOf':
            return SelectValueRow({currentVal: value.val, name: param.name, allowedVals: param.limits, updateVal, errMsg})
        case 'less':
            min = Number.MIN_VALUE
            max = param.constraint.limits[0]-Number.EPSILON
            step = 0.01

            return EnterNumberRow({name, currentVal, step, max, updateVal})
        case 'lessEqual':
            min = Number.MIN_VALUE
            max = param.constraint.limits[0]
            step = 0.01

            return EnterNumberRow({name, currentVal, step, max, updateVal})
        case 'greater':
            min = param.constraint.limits[0]+Number.EPSILON
            max = Number.MAX_VALUE
            step = 0.01

            return EnterNumberRow({name, currentVal, step, min, updateVal})
        case 'greaterEqual':
            min = param.constraint.limits[0]
            max = Number.MAX_VALUE
            step = 0.01

            return EnterNumberRow({name, currentVal, step, min, updateVal})
        case 'rangeIncl':
            min = param.constraint.limits[0]
            max = param.constraint.limits[1]
            step = computeStep(param.limits[0], param.limits[1])

            return EnterNumberRow({name, currentVal, step, min, max, updateVal})
        case 'rangeExcl':
            min = param.constraint.limits[0]
            max = param.constraint.limits[1]-Number.EPSILON
            step = computeStep(param.limits[0], param.limits[1])

            return EnterNumberRow({name, currentVal, step, min, max, updateVal})
    }

    console.log(`unsupported float64 constraint type ${param.constraint.type}`)

    return null
}

const ParamRow = ({param, value, errMsg}) => {
    switch (param.type) {
    case "int":
        return IntParamRow({param, errMsg, value})
    case "float64":
        return FloatParamRow({param, errMsg, value})
    }
    
    console.log(`unsupported param type ${param.type}`)

    return null;
}

function validateParamVal(param, value) {
    switch (param.type) {
    case "IntEnum":
    case "FltEnum":
    case "StrEnum":
        if (param.limits.indexOf(value) == -1) {
            return new Error(`invalid value ${value} for param ${param.name}: not one of enum values ${param.limits}`)
        }
        break
    case "IntRange":
    case "FltRange":
        if (value < param.limits[0] || value > param.limits[1]) {
            return new Error(`invalid value ${value} for param ${param.name}: not in range [${param.limits[0]}, ${param.limits[1]}]`)
        }
        break
    }
        
    return null;
}

const EditParamValsModal = ({params, paramVals, onComplete}) => {
    const closed = van.state(false)

    params.sort((a,b) => {
        if (a.name < b.name) {
            return -1
        }

        if (a.name > b.name) {
            return 1
        }

        return 0
    })

    const paramValsWorking = Object.fromEntries(params.map(p => {
        const nonDefaultVal = paramVals[p.name]

        return [p.name, van.state(nonDefaultVal ?? p.default)]
    }))
    const errMessages = Object.fromEntries(params.map(p => {
        return [p.name, van.state("")]
    }))

    const paramValTableBody = tbody({class:"table-auto"}); 
    const paramValTable = Table({
        columnNames: ["Name", "Value", ""],
        tableBody: paramValTableBody,
    })
    const rows = params.map(p => {
        return ParamRow({
            param: p,
            value: paramValsWorking[p.name],
            errMsg: errMessages[p.name],
        })
    })
    
    van.add(paramValTableBody, rows)

    const cancelBtn = ButtonCancel({
        child: "Cancel",
        onclick: () => closed.val = true,
    })
    const okBtn = Button({
        child: "OK",
        disabled: van.derive(() => {
            return Object.values(errMessages).map(msg => msg.val.length > 0).reduce((result, current) => result || current, false)
        }),
        onclick: () => {
            const nonDefaultVals = {}

            Object.entries(paramValsWorking).forEach(([name, value]) => {
                if (value.val !== params.find(p => p.name === name).default) {
                    nonDefaultVals[name] = value.val
                }
            })

            onComplete(nonDefaultVals)

            closed.val = true
        },
    })
    const buttons = ButtonGroup({buttons: [cancelBtn, okBtn], moreClass: "self-end"})
    const modal = ModalBackground(
        div(
            {id: "foreground", class: "flex flex-col block p-16 rounded-lg bg-white min-w-[25%] max-w-[50%]"},
            p({class: "text-lg font-medium font-bold text-center"}, "Edit Parameters"),
            paramValTable,
            buttons,
        )
    )

    van.add(document.body, () => closed.val ? null : modal);
}

export {EditParamValsModal, validateParamVal};