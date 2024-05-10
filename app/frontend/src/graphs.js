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

const createGraph = async (item) => {
    console.log("creating graph", item);

    const resp = await Post({route: '/graphs', content: item});

    if (resp.status != 204) {
        console.log("failed to create graph", await resp.json());

        return false
    }

    console.log(`created graph %s`, item.id);

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

const GraphCard = ({id, name}) => {
    const deleted = van.state(false);
    const viewBtn = ButtonAct({
        text: "",
        onclick: () => routeTo('graphs', [id]),
    });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => deleted.val = true,
    });

    viewBtn.classList.add("fa-regular");
    viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    return () => deleted.val ? null : div(
        {class: "block rounded-lg p-6 border h-250 w-250"},
        p({class: "text-lg font-medium font-bold text-center mb-6"}, name),
        viewBtn,
        deleteBtn,
    )
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
    const name = van.state(RandomName())

    return div(
        {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-200 rounded-md"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Name"),
        div(
            div(
                label({for: "name"}, "Name"),
                input({
                    id: "name",
                    class: "block px-5 py-5 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring",
                    type: "text",
                    value: name,
                    placeholder: "Unique, non-empty name",
                    oninput: e => name.val = e.target.value,
                }),
            ),
        ),
        div(
            {class:"mt-4 flex justify-center"},
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
    const cardsArea = div(
        {class:"grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-x-10 gap-y-10"},
    )
    
    getGraphs().then(
        (items) => {
            const cards = items.map(item => GraphCard({id: item.id, name: item.name}));

            van.add(cardsArea, cards);
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

                            createGraph(graphItem).then((ok) => {
                                if (ok) {
                                    van.add(cardsArea, GraphCard({id: id, name: name}));
                                    
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

    return van.add(cardsArea, addGraphBtn);
}

export default Graphs;