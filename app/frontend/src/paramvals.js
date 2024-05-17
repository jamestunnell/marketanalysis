import van from "vanjs-core"

import capitalize from './capitalize.js';

const { div, input, label, li, option, select} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
const minSteps = 100.0

function computeStep(min, max) {
    const ratio = max / min;
    
    return (ratio / minSteps) / Math.pow(10.0, Math.ceil(Math.log10(ratio)))
}

const RangeParamVal = ({paramDef, values, step}) => {
    let value = values[paramDef.name]
    if (!value) {
        value = van.state(paramDef.default);
    }

    const min = paramDef.limits[0]
    const max = paramDef.limits[1]
    const labelText = `${capitalize(paramDef.name)}: (${min}-${max}):`;

    return div(
        {class: "flex flex-col"},
        label({for: paramDef.name}, labelText),
        input({
            id: paramDef.name,
            type: "number",
            class: inputClass,
            value: value.val,
            min: min,
            max: max,
            step: step,
            onchange: e => value.val = e.target.value,
        }),
    )
}

const IntRangeParamVal = ({paramDef, values}) => {
    return RangeParamVal({paramDef, values, step: 1})
}

const FltRangeParamVal = ({paramDef, values}) => {
    const step = computeStep(paramDef.limits[0], paramDef.limits[1])

    return RangeParamVal({paramDef, values, step})
}

const EnumParamVal = ({paramDef, values}) => {
    let value = values[paramDef.name]
    if (!value) {
        value = van.state(paramDef.default);
    }

    const options = paramDef.limits.map(allowedStr => {
        let props = {value: allowedStr};
        
        if (allowedStr === value.val) {
            props.selected = "selected";
        }

        return option(props, allowedStr);
    });

    return li(
        label({for: paramDef.name}, capitalize(paramDef.name)),
        select({
            id: paramDef.name,
            class: inputClass,
            oninput: e => value.val = e.target.value,
        }, options),
    )
}

const ParamValItem = (paramDef, values) => {
    switch (paramDef.type) {
    case "IntEnum":
    case "FltEnum":
    case "StrEnum":
        return EnumParamVal({paramDef, values})
    case "IntRange":
        return IntRangeParamVal({paramDef, values})
    case "FltRange":
        return FltRangeParamVal({paramDef, values})
    }
    
    console.log(`unknown param type ${paramDef.type}`)

    return null;
}

function validateParamVal(paramDef, value) {
    switch (paramDef.type) {
    case "IntEnum":
    case "FltEnum":
    case "StrEnum":
        if (paramDef.limits.indexOf(value) == -1) {
            return new Error(`invalid value ${value} for param ${paramDef.name}: not one of enum values ${paramDef.limits}`)
        }
    case "IntRange":
    case "FltRange":
        if (value < paramDef.limits[0] || value > paramDef.limits[1]) {
            return new Error(`invalid value ${value} for param ${paramDef.name}: not in range [${paramDef.limits[0]}, ${paramDef.limits[1]}]`)
        }
    }
        
    return null;
}

export {ParamValItem, validateParamVal};