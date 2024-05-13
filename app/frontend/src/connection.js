import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import {ButtonAct, ButtonCancel} from './buttons.js';

const {div, input, label, p} = van.tags

const ConnectionForm = ({source, target, onOK, onCancel}) => {
    const inputClass = "block px-5 py-5 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    return div(
        {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-300 rounded-md"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Connection"),
        div(
            div(
                label({for: "source"}, "Source Address"),
                input({
                    id: "source",
                    class: inputClass,
                    type: "text",
                    value: source,
                    placeholder: "<block.output>",
                    oninput: e => source.val = e.target.value,
                }),
                label({for: "target"}, "Target Address"),
                input({
                    id: "target",
                    class: inputClass,
                    type: "text",
                    value: target,
                    placeholder: "<block.input>",
                    oninput: e => target.val = e.target.value,
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

const DoConnectionModal = ({connection, handleResult}) => {
    const closed = van.state(false);
    const source = van.state(connection.source);
    const target = van.state(connection.target);

    van.add(
        document.body,
        Modal({closed},
            ConnectionForm({
                source: source,
                target: target,
                onOK: ()=> {
                    console.log("pressed OK")

                    handleResult({source: source.val, target: target.val});

                    console.log("closing modal")

                    closed.val = true;
                },
                onCancel: () => {
                    closed.val = true;
                }
            }),
        ),
    );
}

export {DoConnectionModal};