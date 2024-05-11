import van from "vanjs-core"
import { routeTo } from 'vanjs-router'
import {Modal} from "vanjs-ui"
import { v4 as uuidv4 } from 'uuid';
import { uniqueNamesGenerator, adjectives, colors, animals } from 'unique-names-generator';

import {Table} from './table.js'
import {ButtonAct, ButtonCancel} from './buttons.js'
import {Delete, Get, Post} from './backend.js'

const {div, input, label, p, tbody, td, tr} = van.tags

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

    // Avoid Fetch failed loading
    await resp.text();

    console.log(`created graph %s`, item.id);

    return true;
}

const deleteGraph = async (id) => {
    console.log("deleting graph %s", id);

    const resp = await Delete(`/graphs/${id}`);

    if (resp.status != 204) {
        console.log("failed to delete graph", await resp.json());

        return false
    }

    // Avoid Fetch failed loading
    await resp.text();

    console.log(`deleted graph %s`, id);

    return true;
}

const truncateString = (id, len) => {
    if (id.length > len) {
        return id.substring(0, len) + "..."
    }
    
    return id
}

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

const GraphTableRow = ({id, name}) => {
    const deleted = van.state(false)

    const viewBtn = ButtonAct({
        text: "",
        onclick: () => routeTo('graphs', [id]),
    });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => {
            deleteGraph(id).then(ok => {
                if (ok) {
                    deleted.val = true
                }
            })
        },
    });

    viewBtn.classList.add("fa-regular");
    viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    return () => deleted.val ? null : tr(
        {class: "border border-solid"},
        td({class: "px-6 py-4"}, truncateString(id, 8)),
        td({class: "px-6 py-4"}, name),
        td(
            {class: "px-6 py-4"},
            div({class:"flex flex-row"}, viewBtn, deleteBtn)
        ),
    )
}

const Graphs = () => {
    const columnNames = ["ID", "Name", ""]
    const tableBody = tbody({class:"table-auto"});

    getGraphs().then(
        (items) => {
            const rows = items.map(item => GraphTableRow({id: item.id, name: item.name}));

            van.add(tableBody, rows);
        }
    );

    const addGraphBtn = ButtonAct({
        text: "Add New",
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
                                    van.add(tableBody, GraphTableRow({id: id, name: name, }));
                                    
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

    return div(
        addGraphBtn,
        Table({columnNames: columnNames, tableBody: tableBody}),
    )
}

export default Graphs;