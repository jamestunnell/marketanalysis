import van from "vanjs-core"

import capitalize from './capitalize.js';

const { div, input, label, li, option, select} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
const minSteps = 100.0

function computeStep(min, max) {
    const ratio = max / min;
    
    return (ratio / minSteps) / Math.pow(10.0, Math.ceil(Math.log10(ratio)))
}

const RangeParamVal = ({param, value, step, onChange}) => {
    const min = param.limits[0]
    const max = param.limits[1]
    const labelText = `${capitalize(param.name)}: (${min}-${max}):`;

    return div(
        {class: "flex flex-col"},
        label({for: param.name}, labelText),
        input({
            id: param.name,
            type: "number",
            class: inputClass,
            value: value.val,
            min: min,
            max: max,
            step: step,
            onchange: e => {
                const newVal = Number(e.target.value)

                value.val = newVal

                onChange(newVal)
            },
        }),
    )
}

const IntRangeParamVal = ({param, value, onChange}) => {
    return RangeParamVal({param, value, step: 1, onChange})
}

const FltRangeParamVal = ({param, value, onChange}) => {
    const step = computeStep(param.limits[0], param.limits[1])

    return RangeParamVal({param, value, step, onChange})
}

const EnumParamVal = ({param, currentVal, updateVal}) => {
    const options = param.limits.map(allowedVal => {
        let props = {value: allowedVal};
        
        if (allowedVal === currentVal) {
            props.selected = "selected";
        }

        return option(props, allowedVal);
    });

    return li(
        label({for: param.name}, capitalize(param.name)),
        select({
            id: param.name,
            class: inputClass,
            oninput: e => updateVal(e.target.value),
        }, options),
    )
}

const ParamValItem = ({param, value, onChange}) => {
    switch (param.type) {
    case "IntEnum":
        return EnumParamVal({
            param,
            currentVal: value.val,
            updateVal: (strVal) => {
                const newVal = parseInt(strVal)
                value.val = newVal

                onChange(newVal)
            },
        })
    case "FltEnum":
        return EnumParamVal({
            param,
            currentVal: value.val,
            updateVal: (strVal) => {
                const newVal = parseFloat(strVal)
                value.val = newVal

                onChange(newVal)
            }
        })
    case "StrEnum":
        return EnumParamVal({
            param,
            currentVal: value.val,
            updateVal: (strVal) => {
                value.val = strVal

                onChange(strVal)
            },
        })
    case "IntRange":
        return IntRangeParamVal({param, value, onChange})
    case "FltRange":
        return FltRangeParamVal({param, value, onChange})
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
    case "IntRange":
    case "FltRange":
        if (value < param.limits[0] || value > param.limits[1]) {
            return new Error(`invalid value ${value} for param ${param.name}: not in range [${param.limits[0]}, ${param.limits[1]}]`)
        }
    }
        
    return null;
}

export {ParamValItem, validateParamVal};