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

const RangeParamRow = ({param, currentVal, step, updateVal, errMsg}) => {
    const min = param.limits[0]
    const max = param.limits[1]
    const labelText = `${capitalize(param.name)}: (${min}-${max}):`;

    const status = ButtonIconTooltip({
        icon: () => (errMsg.val.length > 0) ? IconError() : IconCheck(),
        tooltipText: van.derive(() => (errMsg.val.length > 0) ? `Value is invalid: ${errMsg.val}` : "Value is valid"),
    });
    const rowItems = [
        label({for: param.name}, labelText),
        input({
            id: param.name,
            type: "number",
            class: inputClass,
            value: currentVal,
            min: min,
            max: max,
            step: step,
            onchange: e => {
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

const IntRangeParamRow = ({param, currentVal, updateVal, errMsg}) => {
    return RangeParamRow({param, currentVal, step: 1, updateVal, errMsg})
}

const FltRangeParamRow = ({param, currentVal, updateVal, errMsg}) => {
    const step = computeStep(param.limits[0], param.limits[1])

    return RangeParamRow({param, currentVal, step, updateVal, errMsg})
}

const EnumParamRow = ({param, currentVal, updateVal, errMsg}) => {
    const options = param.limits.map(allowedVal => {
        let props = {value: allowedVal};
        
        if (allowedVal === currentVal) {
            props.selected = "selected";
        }

        return option(props, allowedVal);
    });
    const status = ButtonIconTooltip({
        icon: () => (errMsg.val.length > 0) ? IconError() : IconCheck(),
        tooltipText: van.derive(() => (errMsg.val.length > 0) ? `Value is invalid: ${errMsg.val}` : "Value is valid"),
    });

    const rowItems = [
        label({for: param.name}, capitalize(param.name)),
        select({
            id: param.name,
            class: inputClass,
            oninput: e => {
                const err = updateVal(e.target.value)
                if (err) {
                    errMsg.val = err.message
                } else {
                    errMsg.val = ""
                }
            },
        }, options),
        status,
    ]

    return TableRow(rowItems)
}

const ParamRow = ({param, value, errMsg}) => {
    switch (param.type) {
    case "IntEnum":
        return EnumParamRow({
            param,
            errMsg,
            currentVal: value.val,
            updateVal: (strVal) => {
                const newVal = parseInt(strVal)
                if (isNaN(newVal)) {
                    return new Error(`${strVal} is not an integer`)
                }

                value.val = newVal

                return validateParamVal(param, newVal)
            },
        })
    case "FltEnum":
        return EnumParamRow({
            param,
            errMsg,
            currentVal: value.val,
            updateVal: (strVal) => {
                const newVal = parseFloat(strVal)
                if (isNaN(newVal)) {
                    return new Error(`${strVal} is not an float`)
                }

                value.val = newVal

                return validateParamVal(param, newVal)
            }
        })
    case "StrEnum":
        return EnumParamRow({
            param,
            errMsg,
            currentVal: value.val,
            updateVal: (strVal) => {
                value.val = strVal
                
                return validateParamVal(param, strVal)
            },
        })
    case "IntRange":
        return IntRangeParamRow({
            param,
            errMsg,
            currentVal: value.val,
            updateVal: (strVal) => {
                const newVal = parseInt(strVal)
                if (isNaN(newVal)) {
                    return new Error(`${strVal} is not an integer`)
                }

                value.val = newVal

                return validateParamVal(param, value.val)
            },
        })
    case "FltRange":
        return FltRangeParamRow({
            param,
            errMsg,
            currentVal: value.val,
            updateVal: (strVal) => {
                const newVal = parseFloat(strVal)
                if (isNaN(newVal)) {
                    return new Error(`${strVal} is not an float`)
                }

                value.val = newVal

                return validateParamVal(param, value.val)
            }
        })
    }
    
    console.log(`unknown param type ${param.type}`)

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