import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import {Button, ButtonCancel} from './buttons.js';

const {div, input, label, p} = van.tags

const ConnectionForm = ({source, target, onOK, onCancel}) => {
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    return div(
        {class: "flex flex-col rounded-md space-y-4"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Connection"),
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
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({onclick: onCancel, child: "Cancel"}),
            Button({onclick: onOK, child: "OK"}),
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