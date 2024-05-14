import van from "vanjs-core"
import { routeTo } from 'vanjs-router'
import {Modal} from "vanjs-ui"
import { v4 as uuidv4 } from 'uuid';
import { uniqueNamesGenerator, adjectives, animals } from 'unique-names-generator';

import {Delete, Get, Post} from './backend.js'
import { Button, ButtonCancel, ButtonDanger } from "./buttons.js";
import {IconDelete, IconView} from './icons.js'
import {Table, TableRow} from './table.js'

const {div, input, label, p, tbody} = van.tags

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
    return uniqueNamesGenerator({ dictionaries: [adjectives, animals] });
}

const GraphNameForm = ({onOK, onCancel}) => {
    const name = van.state(RandomName())

    return div(
        {class: "flex flex-col rounded-md space-y-4"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Name"),
        label({for: "name"}, "Name"),
        input({
            id: "name",
            class: "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring",
            type: "text",
            value: name,
            placeholder: "Unique, non-empty name",
            oninput: e => name.val = e.target.value,
        }),
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({child: "Cancel", onclick: onCancel}),
            Button({child: "OK", onclick: ()=> onOK({name: name.val})}),
        ),
    )
}

const GraphTableRow = ({id, name}) => {
    const deleted = van.state(false)

    const viewBtn = Button({
        child: IconView(),
        onclick: () => routeTo('graphs', [id]),
    });
    const deleteBtn = ButtonDanger({
        child: IconDelete(),
        onclick: () => {
            deleteGraph(id).then(ok => {
                if (ok) {
                    deleted.val = true
                }
            })
        },
    });

    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [name, truncateString(id, 8), buttons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const Graphs = () => {
    const columnNames = ["Name", "ID", ""]
    const tableBody = tbody({class:"table-auto"});

    getGraphs().then((graphs) => {
        const rows = graphs.map(g => GraphTableRow({id: g.id, name: g.name}));

        van.add(tableBody, rows);
    });

    const newGraphBtn = Button({
        child: "New Graph",
        onclick: () => {
            const closed = van.state(false)

            van.add(
                document.body,
                Modal({closed},
                    GraphNameForm({
                        onOK: ({name})=> {
                            const id  = uuidv4();
                            const graphItem = {id: id, name: name, blocks: [], connections: []};

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
        div(
            {class: "container flex flex-row-reverse p-4"},
            newGraphBtn,
        ),
        Table({columnNames: columnNames, tableBody: tableBody}),
    )
}

export default Graphs;