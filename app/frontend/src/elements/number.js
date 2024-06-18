import van from "vanjs-core"

import {INPUT_CLASS} from './input.js'

const {input} = van.tags

const NumberRange = ({min=null, max=null, id="", parse, step, validate, value}) => {
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

            validate(newVal)
        },
    })
}

const NumberEnum = ({id="", parse, allowed, validate, value}) => {
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

            validate(newVal)
        },
    }, options)
}

const IntRange = ({constraint, id="", error, value}) => {
    return NumberRange({
        id, value,
        step: 1,
        min: constraint.getMin(),
        max: constraint.getMax(),
        parse: parseInt,
        validate: (val) => {
            error.val = constraint.validate(val)
        },
    })
}

const IntEnum = ({constraint, id="", error, value}) => {
    return NumberEnum({
        id, value,
        allowed: constraint.getAllowed(),
        parse: parseInt,
        validate: (val) => {
            error.val = constraint.validate(val)
        },
    })
}

const FloatRange = ({constraint, id="", error, value}) => {
    return NumberRange({
        id, value,
        step: 0.01,
        min: constraint.getMin(),
        max: constraint.getMax(),
        parse: parseFloat,
        validate: (val) => {
            error.val = constraint.validate(val)
        },
    })
}

const FloatEnum = ({constraint, id="", error, value}) => {
    return NumberEnum({
        id, value,
        allowed: constraint.getAllowed(),
        parse: parseFloat,
        validate: (val) => {
            error.val = constraint.validate(val)
        },
    })
}

export {IntRange, IntEnum, FloatRange, FloatEnum}