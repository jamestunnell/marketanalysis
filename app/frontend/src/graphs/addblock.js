import van from "vanjs-core"

import { Button, ButtonCancel } from '../elements/buttons.js';
import { ModalBackground, ModalForeground } from "../modal.js";
import { truncateString } from "../truncatestring.js";

const {div, option, p, select, span} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-gray-500 focus:outline-none focus:ring";

const AddBlockForm = ({infoByType, blockNames, onComplete, onCancel}) => {
    const types = Object.keys(infoByType)
    const selectedType = van.state(types[0])
    const description = van.derive(() => infoByType[selectedType.val].description)
    const options = types.map((t) => {
        return option({value: t, selected: (t === selectedType.val)}, t);
    })
    const selectType = select(
        {
            id: "type",
            class: inputClass,
            onchange: (e) => selectedType.val = e.target.value,
        },
        options,
    )
    const ok = Button({
        child: "OK",
        onclick: () => {
            const info = infoByType[selectedType.val]
            let name = truncateString(selectedType.val, 3).toLowerCase()

            if (blockNames.indexOf(name) >= 0) {
                let i = 2
                const candidate = () => {return `${name}${i}`}

                while(blockNames.indexOf(candidate()) >= 0) {
                    i++
                }

                name = candidate()
            }

            const config = {type: selectedType.val, name, parameters: [], inputs: [], outputs: []}
    
            onComplete({info, config})
        },
    })
    const cancel = ButtonCancel({child: "Cancel", onclick: onCancel})

    return div(
        {class: "flex flex-col max-w-250"},
        div(
            {class: "grid grid-cols-2"},
            span({class: "min-w-24 max-w-24"}, p({class: "text-md font-medium font-bold"}, "Type")),
            span({class: "min-w-48 max-w-48"}, selectType),
            span({class: "min-w-24 max-w-24"}, p({class: "text-md font-medium font-bold"}, "Description")),
            span({class: "min-w-48 max-w-48"}, description)
        ),
        div({class:"mt-4 flex flew-row-reverse"}, ok, cancel),
    )
}

const AddBlockModal = ({infoByType, blockNames, handleResult}) => {
    const closed = van.state(false)
    const onComplete = ({config, info}) => {
        handleResult({config, info})

        closed.val = true
    }
    const onCancel = () => closed.val = true;
    const form = AddBlockForm({infoByType, blockNames, onComplete, onCancel})
    const modal = ModalBackground(ModalForeground({}, form))

    van.add(document.body, () => closed.val ? null : modal);
}

export {AddBlockForm, AddBlockModal}