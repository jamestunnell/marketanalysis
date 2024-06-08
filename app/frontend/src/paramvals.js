import van from "vanjs-core"

import capitalize from './capitalize.js';
import { ModalBackground } from "./modal.js";
import { Button, ButtonCancel, ButtonIconTooltip } from "./buttons.js";
import { ButtonGroup } from "./buttongroup.js";
import { IconCheck, IconError } from "./icons.js";
import { Table, TableRow } from "./table.js";

const { div, input, label, option, p, select, tbody} = van.tags

const inputClass = "block border border-gray-200 rounded-md focus:border-gray-500 focus:outline-none focus:ring";

const EnterNumberRow = ({currentVal, param, step, updateVal}) => {
    const errMsg = van.state("")
    const inputProps = {
        id: param.name,
        type: "number",
        class: inputClass,
        value: currentVal,
        step,
        onchange: (e) => {
            const err = updateVal(e.target.value)
            if (err) {
                errMsg.val = err.message
            } else {
                errMsg.val = ""
            }
        },
    }

    let constraintText
    
    switch (param.constraint.type) {
        case 'less':
            inputProps.max = param.constraint.limits[0] - step
            constraintText = `< ${param.constraint.limits[0]}`
            break
        case 'lessEqual':
            inputProps.max = param.constraint.limits[0]
            constraintText = `<= ${param.constraint.limits[0]}`
            break
        case 'greater':
            inputProps.min = param.constraint.limits[0] + step
            constraintText = `> ${param.constraint.limits[0]}`
            break
        case 'greaterEqual':
            inputProps.min = param.constraint.limits[0]
            constraintText = `>= ${param.constraint.limits[0]}`
            break
        case 'rangeIncl':
            inputProps.min = param.constraint.limits[0]
            inputProps.max = param.constraint.limits[1]
            constraintText = `[${param.constraint.limits[0]}, ${param.constraint.limits[1]}]`
            break
        case 'rangeExcl':
            inputProps.min = param.constraint.limits[0]
            inputProps.max = param.constraint.limits[1] - step
            constraintText = `[${param.constraint.limits[0]}, ${param.constraint.limits[1]})`
            break
    }

    const status = ButtonIconTooltip({
        icon: () => (errMsg.val.length > 0) ? IconError() : IconCheck(),
        tooltipText: van.derive(() => (errMsg.val.length > 0) ? `Value is invalid: ${errMsg.val}` : "Value is valid"),
    });
    const rowItems = [
        label({for: param.name}, param.name),
        constraintText,
        input(inputProps),
        status,
    ]

    return TableRow(rowItems)
}

const SelectValueRow = ({param, currentVal, updateVal}) => {
    const errMsg = van.state("")
    const options = param.constraint.limits.map(allowedVal => {
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
    const constraintText = ""

    const rowItems = [
        label({for: param.name}, capitalize(param.name)),
        constraintText,
        select({
            id: param.name,
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
    
        return validateParamVal({param, value: newVal})
    }

    switch (param.constraint.type) {
        case 'oneOf':
            return SelectValueRow({currentVal: value.val, param, updateVal})
        case 'less':
        case 'lessEqual':
        case 'greater':
        case 'greaterEqual':
        case 'rangeIncl':
        case 'rangeExcl':
        case 'none':
            return EnterNumberRow({currentVal: value.val, step:1, param, updateVal})
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

        return validateParamVal({param, value: newVal})
    }

    switch (param.constraint.type) {
        case 'oneOf':
            return SelectValueRow({currentVal: value.val, param, updateVal})
        case 'less':
        case 'lessEqual':
        case 'greater':
        case 'greaterEqual':
        case 'rangeIncl':
        case 'rangeExcl':
        case 'none':
            return EnterNumberRow({currentVal: value.val, step:0.01, param, updateVal})
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

function validateParamVal({param, value}) {
    const limits = param.constraint.limits

    let err

    switch (param.constraint.type) {
    case "oneOf":
        if (limits.indexOf(value) === -1) {
            err = new Error(`${value} is not one of ${limits}`) 
        }
        break
    
    case "less":
        if (value >= limits[0]) {
            err = new Error(`${value} is not < ${limits[0]}`) 
        }
        break
    case "lessEqual":
        if (value > limits[0]) {
            err = new Error(`${value} is not <= ${limits[0]}`) 
        }
        break
    case "greater":
        if (value <= limits[0]) {
            err = new Error(`${value} is not > ${limits[0]}`) 
        }
        break
    case "greaterEqual":
        if (value < limits[0]) {
            err = new Error(`${value} is not >= ${limits[0]}`) 
        }
        break
    case "rangeIncl":
        if (value < limits[0] || value > limits[0]) {
            err = new Error(`${value} is not in range [${limits[0]}, ${limits[1]}]`) 
        }
        break
    case "rangeExcl":
        if (value < limits[0] || value >= limits[0]) {
            err = new Error(`${value} is not in range [${limits[0]}, ${limits[1]})`) 
        }
        break
    }
        
    return err;
}

const ParamValsTable = ({params, paramVals}) => {
    const names = params.map(p => p.name).sort()
    const rows = names.map(name => {
        return ParamRow({
            param: params.find(p => p.name === name),
            value: paramVals[name],
        })
    })

    return Table({
        columnNames: ["Name", "Constraint", "Value", ""],
        tableBody: tbody({class:"table-auto"}, rows),
    })
}

export {ParamValsTable, validateParamVal};