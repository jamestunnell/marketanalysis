import van from "vanjs-core"

import {INPUT_CLASS} from './input.js'

const {input} = van.tags

const Textbox = ({id="", value=van.state(""), placeholder=""}) => {
    return input({
        id,
        value,
        placeholder,
        class: INPUT_CLASS,
        type: "text",
        onchange: e => text.val = e.target.value,
    })
}

export default Textbox