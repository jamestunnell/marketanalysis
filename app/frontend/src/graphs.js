import van from "vanjs-core"
import { routeTo } from 'vanjs-router'
import {Modal} from "vanjs-ui"
import { v4 as uuidv4 } from 'uuid';
import { uniqueNamesGenerator, adjectives, colors, animals } from 'unique-names-generator';

import {ButtonAct, ButtonCancel} from './buttons.js'
import {Get, Post} from './backend.js'

const {button, div, h2, input, label, p} = van.tags

const getGraphs = async () => {
    console.log("getting graphs");

    const resp = await Get('/graphs');

    if (resp.status != 200) {
        console.log("failed to get graphs", await resp.json());

        return []
    }

    const d = await resp.json();

    console.log(`received ${d.graphs.length} graphs`, d.graphs);

    return d.graphs;
}

const addGraph = async (item) => {
    console.log("adding graph", item);

    const resp = await Post({route: '/graphs', content: item});

    if (resp.status != 204) {
        console.log("failed to add graph", await resp.json());

        return false
    }

    console.log(`added graph %s`, item.id);

    return true;
}

// const delGraph = async (symbol) => {
//     console.log("deleting graph");

//     const resp = await fetch(`${BASE_URL}/graphs/${symbol}`, {
//         method: 'DELETE',
//         credentials: 'same-origin'
//     });

//     console.log('delete graph result:', resp.status)

//     return resp.status === 204 
// }


const Btn = ({onclick}) => {
    return button(
        {
            class: "block rounded-lg p-6 border h-100 w-100",
            onclick: onclick,
        },
    );
}

const GraphBtn = ({id, name}) => {
    console.log(`making graph button {id: ${id}, name: ${name}}`);

    // const deleted = van.state(false);
    const btn = Btn({onclick: () => routeTo('graphs', [id])});

    return van.add(btn, h2(name));
    // const editBtn = ButtonAct({
    //     text: "",
    //     onclick: () => 
    // });
    // const deleteBtn = ButtonAct({
    //     text: "",
    //     onclick: () => deleted.val = true,
    // });

    // editBtn.classList.add("fa-solid");
    // editBtn.classList.add("fa-pen-to-square");

    // deleteBtn.classList.add("fa-solid");
    // deleteBtn.classList.add("fa-trash");


}

// const ID_PREVIEW_LEN = 8;

// const truncateString = (id, len) => {
//     if (id.length > len) {
//         return id.substring(0, len) + "..."
//     }
    
//     return id
// }

const RandomName = () => {
    return uniqueNamesGenerator({ dictionaries: [adjectives, colors, animals] });
}

const GraphNameForm = ({onOK, onCancel}) => {
    const editBoxClass = "block w-full px-4 py-2 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
    const name = van.state(RandomName())

    return div(
        {class: "w-200 space-y-6"},
        p({class: "text-lg font-medium font-bold"}, "Graph Name"),
        div(
            {class: "grid grid-cols-1 gap-6 mt-4"},
            div(
                label({for: "name"}, "Name"),
                input({id: "name", class: editBoxClass, type: "text", value: name, oninput: e => name.val = e.target.value, placeholder: "Unique, non-empty name"}),
            ),
        ),
        div(
            {class:"mt-4 flex justify-end"},
            ButtonCancel({text: "Cancel", onclick: () => onCancel()}),
            ButtonAct({
                text: "OK",
                onclick: async () => {
                    onOK({name: name.val})
                },
            }),
        ),
    )
}

const Graphs = () => {
    const graphsArea = div(
        {class:"flex flex-wrap overflow-y-scroll px-6 py-4"},
    )
    
    getGraphs().then(
        (items) => {
            const btns = items.map(item => GraphBtn({id: item.id, name: item.name}));

            van.add(graphsArea, btns);
        }
    );

    const addGraphBtn = Btn({
        onclick: () => {
            const closed = van.state(false)

            van.add(
                document.body,
                Modal({closed},
                    GraphNameForm({
                        onOK: ({name})=> {
                            const id  = uuidv4();
                            const graphItem = {id: id, name: name, blocks: {}, connections: []};

                            addGraph(graphItem).then((ok) => {
                                if (ok) {
                                    van.add(graphsArea, GraphBtn({id: id, name: name}));
                                    
                                    closed.val = true;
                                }
                            });
                        },
                        onCancel: () => {
                            closed.val = true
                        }
                    }),
                ),
            );
        },
    });

    addGraphBtn.classList.add("fa-solid");
    addGraphBtn.classList.add("fa-plus");
    addGraphBtn.classList.add("order-last");

    return van.add(graphsArea, addGraphBtn);
}

export default Graphs;