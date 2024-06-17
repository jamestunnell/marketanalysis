import van from "vanjs-core"

import {INPUT_CLASS} from './input.js'

const {input} = van.tags

const Checkbox = ({id="", checked=van.state(false)}) => {
    return input({
        id, checked,
        class: INPUT_CLASS,
        type: "checkbox",
        onchange: e => checked.val = e.target.checked,
    })
}

export default Checkbox