import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import { Get } from './backend.js'
import {Button, ButtonCancel} from './buttons.js';

const {div, input, label, option, p, select} = van.tags

const getBlockInfos = async () => {
    console.log("getting block infos");

    const resp = await Get('/blocks');

    if (resp.status != 200) {
        console.log("failed to get block infos", await resp.json());

        return []
    }

    const d = await resp.json();

    console.log(`received ${d.blocks.length} block infos`, d.blocks);

    return d.blocks;
}

const BlockForm = ({name, type, onOK, onCancel}) => {
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    const typeSelect = select({
        id: "type",
        class: inputClass,
        oninput: (e) => type.val = e.target.value,
    });
    
    getBlockInfos().then(blockInfos => {
        van.add(typeSelect, blockInfos.map(info => {
            const t = info.type;
            
            let props = {value: t}
            
            if (t === type.val) {
                props.selected = "selected";
            }

            return option(props, t)
        }))
    });

    return div(
        {class: "flex flex-col rounded-md space-y-4"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Block"),
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
        typeSelect,
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({onclick: onCancel, child: "Cancel"}),
            Button({onclick: onOK, child: "OK"}),
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