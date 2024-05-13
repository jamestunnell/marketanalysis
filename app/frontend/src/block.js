import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import {ButtonAct, ButtonCancel} from './buttons.js';

const {div, input, label, p} = van.tags

const BlockForm = ({name, type, onOK, onCancel}) => {
    const inputClass = "block px-5 py-5 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    return div(
        {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-300 rounded-md"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Block"),
        div(
            div(
                label({for: "name"}, "Name"),
                input({
                    id: "name",
                    class: inputClass,
                    type: "text",
                    value: name,
                    placeholder: "Non-empty, unique",
                    oninput: e => name.val = e.target.value,
                }),
                label({for: "type"}, "Type"),
                input({
                    id: "type",
                    class: inputClass,
                    type: "text",
                    value: type,
                    placeholder: "Valid block type",
                    oninput: e => type.val = e.target.value,
                }),
            ),
        ),
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({text: "Cancel", onclick: onCancel}),
            ButtonAct({text: "OK", onclick: onOK}),
        ),
    )
}

const DoBlockModal = ({block, handleResult}) => {
    const closed = van.state(false)

    const name = van.state(block.name);
    const type = van.state(block.type);
    const paramVals = van.state(block.paramVals);
    const recording = van.state(block.recording);

    van.add(
        document.body,
        Modal({closed},
            BlockForm({
                name: name,
                type: type,
                // paramVals: paramVals,
                // recording: recording,
                onOK: ()=> {
                    handleResult({
                        name: name.val,
                        type: type.val,
                        paramVals: paramVals.val,
                        recording: recording.val,
                    });

                    closed.val = true;
                },
                onCancel: () => {
                    closed.val = true;
                }
            }),
        ),
    );
}

export {DoBlockModal};