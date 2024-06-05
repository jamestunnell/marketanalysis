import van from "vanjs-core"

import capitalize from './capitalize.js';
import { ModalBackground } from "./modal.js";
import { Button, ButtonCancel, ButtonIconTooltip } from "./buttons.js";
import { ButtonGroup } from "./buttongroup.js";
import { IconCheck, IconError } from "./icons.js";
import { Table, TableRow } from "./table.js";

const { div, input, label, option, p, select, tbody} = van.tags

const inputClass = "block border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

const EnterNumberRow = ({currentVal, param, step, updateVal, errMsg}) => {
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
            return SelectValueRow({currentVal: value.val, param, updateVal, errMsg})
        case 'less':
        case 'lessEqual':
        case 'greater':
        case 'greaterEqual':
        case 'rangeIncl':
        case 'rangeExcl':
        case 'none':
            return EnterNumberRow({currentVal: value.val, step:1, param, updateVal, errMsg})
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
            return SelectValueRow({currentVal: value.val, param, updateVal, errMsg})
        case 'less':
        case 'lessEqual':
        case 'greater':
        case 'greaterEqual':
        case 'rangeIncl':
        case 'rangeExcl':
        case 'none':
            return EnterNumberRow({currentVal: value.val, step:0.01, param, updateVal, errMsg})
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
        columnNames: ["Name", "Constraint", "Value", ""],
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