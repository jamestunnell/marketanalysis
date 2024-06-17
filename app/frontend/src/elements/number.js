import van from "vanjs-core"

import {INPUT_CLASS} from './input.js'
import { constraintMinMax } from "./constraints.js"

const {input} = van.tags

const Integer = ({id="", value=van.state(0)}) => {
    return Number({step: 1, id, value})
}

const Number = ({step, id="", value=van.state(0)}) => {
    const props = {
        id,
        value,
        step,
        class: INPUT_CLASS,
        type: "number",
        onchange: (e) => {
            const newVal = parseInt(strVal)
            if (isNaN(newVal)) {
                return new Error(`${strVal} is not an integer`)
            }
        
            value.val = newVal
        },

        const newVal = parseInt(strVal)
        if (isNaN(newVal)) {
            return new Error(`${strVal} is not an integer`)
        }
    
        value.val = newVal
    }

    constraintMinMax({
        constraint, step,
        onMin: (val) => {props.min = val},
        onMax: (val) => {props.max = val}
    })

    return input(props)
}

const ConstrainedNumber = ()

export {Integer, Float}