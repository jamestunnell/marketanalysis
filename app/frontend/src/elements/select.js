import van from "vanjs-core"

import {INPUT_CLASS} from './input.js'

const {select} = van.tags

const Select = ({onchange, id="", options=[]}) => {
    return select({
        id,
        class: INPUT_CLASS,
        onchange,
    }, options)
}

export default Select