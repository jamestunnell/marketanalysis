import van from "vanjs-core"

import {INPUT_CLASS} from './input.js'

const {input} = van.tags

const NumberRange = ({min=null, max=null, id="", parse, step, value}) => {
    return input({
        id, step, min, max,
        class: INPUT_CLASS,
        type: "number",
        value: value,
        onchange: (e) => {
            const strVal = e.target.value
            const newVal = parse(strVal)
            if (isNaN(newVal)) {
                return new Error(`${strVal} is not an integer`)
            }
        
            value.val = newVal
        },
    })
}

const NumberEnum = ({id="", parse, allowed, value}) => {
    const options = allowed.map(allowedVal => {
        let props = {value: allowedVal};
        
        if (allowedVal === currentVal) {
            props.selected = "selected";
        }

        return option(props, allowedVal);
    })

    return select({
        id,
        class: INPUT_CLASS,
        oninput: (e) => {
            const strVal = e.target.value
            const newVal = parse(strVal)
            if (isNaN(newVal)) {
                return new Error(`${strVal} is not an integer`)
            }
        
            value.val = newVal
        },
    }, options)
}

const IntRange = ({constraint, id="", value}) => {
    const min = constraint.getMin()
    const max = constraint.getMax()
    const step = 1

    return NumberRange({
        id, value, step,
        min: min ? (min.inclusive ? min.value : (min.value + step)): null,
        max: max ? (max.inclusive ? max.value : (max.value - step)): null,
        parse: parseInt,
    })
}

const IntEnum = ({constraint, id="", value}) => {
    return NumberEnum({
        id, value,
        allowed: constraint.getAllowed(),
        parse: parseInt,
    })
}

const FloatRange = ({constraint, id="", value}) => {
    const min = constraint.getMin()
    const max = constraint.getMax()
    const step = 0.01
    
    return NumberRange({
        id, value, step,
        min: min ? (min.inclusive ? min.value : (min.value + step)): null,
        max: max ? (max.inclusive ? max.value : (max.value - step)): null,
        parse: parseFloat,
    })
}

const FloatEnum = ({constraint, id="", value}) => {
    return NumberEnum({
        id, value,
        allowed: constraint.getAllowed(),
        parse: parseFloat,
    })
}

export {IntRange, IntEnum, FloatRange, FloatEnum, NumberRange}