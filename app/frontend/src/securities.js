import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import AppError from "./apperror.js"
import {Delete, Get, Put, Post} from './backend.js'
import {ButtonAct, ButtonCancel} from './buttons.js'

const {button, div, input, option, label, p, select} = van.tags

const TIME_ZONES = Intl.supportedValuesOf('timeZone');
const TZ_NEW_YORK = "America/New_York";

const getSecurities = async () => {
    console.log("getting securities");

    const resp = await Get('/securities');

    const d = await resp.json()

    if (resp.status != 200) {
        console.log("failed to get securities", d);

        return []
    }

    const securities = d.securities;

    console.log("received %d securities", securities.length);

    return securities;
}

const updateSecurity = async (updatedItem) => {
    console.log("updating security", updatedItem);

    const resp = await Put({route: `/securities/${updatedItem.symbol}`, content: updatedItem});
    
    if (resp.status !== 204) {
        const errData = await resp.json();

        console.log('failed to update security', errData);
    
        return errData;
    }

    // Avoid Fetch failed loading
    await resp.text();

    console.log('udpated security', updatedItem);

    return null;
}

const addSecurity = async (newItem) => {
    console.log("adding new security", newItem);

    const resp = await Post({route: '/securities', content: newItem});

    if (resp.status !== 204) {
        const errData = await resp.json();

        console.log('failed to add security', errData);
    
        return errData;
    }

    // Avoid Fetch failed loading
    await resp.text();

    console.log('added security', newItem);

    return null
}

const delSecurity = async (symbol) => {
    console.log("deleting security %s", symbol);

    const resp = await Delete(`/securities/${symbol}`);
                
    if (resp.status !== 204) {
        console.log('failed to delete security', await resp.json());
    
        return false;
    }

    // Avoid Fetch failed loading
    await resp.text();

    console.log('deleted security %s', symbol);

    return true;
}

const SecurityForm = ({title, onCancel, onOK, item}) => {
    const editing = item.symbol !== "";

    const sym = van.state(item.symbol);
    const tz = van.state(item.timeZone);
    const open = van.state(item.open);
    const close = van.state(item.close);
    const errType = van.state("");
    const errMsg = van.state("");
    const errDetails = van.state([]);
    const errHidden = van.state(true)

    const makeItem = () => {
        return {symbol: sym.val, timeZone: tz.val, open: open.val, close: close.val};
    }

    const editBoxClass = "block w-full px-4 py-2 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
    const symbolEditBoxClass = editBoxClass + (editing ? " bg-gray-200" : "");

    return div(
        {class: "space-y-6"},
        p({class: "text-lg font-medium font-bold"}, title),
        div(
            {class: "grid grid-cols-1 gap-6 mt-4"},
            div(
                label({for: "symbol"}, "Symbol"),
                input({id: "symbol", class: symbolEditBoxClass, type: "text", disabled: editing, value: sym, oninput: e => sym.val = e.target.value, placeholder: "SPY, QQQ, etc."}),
            ),
            div(
                label({for: "timeZone"}, "Time Zone"),
                select(
                    {id: "timeZone", class: editBoxClass, oninput: (e) => tz.val = e.target.value},
                    TIME_ZONES.map(x => {
                        let props = {value: x};
                        if (x === item.timeZone) {
                            props.selected = "selected";
                        }
                        
                        return option(props, x)                    
                    }),
                ),    
            ),
            div(
                label({for: "open"}, "Open"),
                input({id: "open", class: editBoxClass, type: "text", value: open, oninput: e => open.val = e.target.value}),
            ),
            div(
                label({for: "close"}, "Close"),
                input({id: "close", class: editBoxClass, type: "text", value: close, oninput: e => close.val = e.target.value}),    
            )
        ),
        AppError({type: errType, msg: errMsg, details: errDetails, hidden: errHidden}),
        div(
            {class:"mt-4 flex justify-end"},
            ButtonCancel({text: "Cancel", onclick: () => onCancel()}),
            ButtonAct({
                text: "OK",
                onclick: async () => {
                    const err = await onOK(makeItem());
                    if (err == null) {
                        return;
                    }
                    
                    errType.val = err.errType;
                    errMsg.val = err.message;
                    errDetails.val = err.details;
                    errHidden.val = false;
                },
            }),
        ),
    )
}

const AddNewButton = ({sidebar, state}) => {
    const btn = ButtonAct({
        text: "",
        onclick: () => {
            const closed = van.state(false)

            van.add(
                document.body,
                Modal(
                    {closed},
                    div({style: "display: flex; justify-content: center;"},
                        SecurityForm({
                            title: "New Security",
                            item: {symbol: "", timeZone: TZ_NEW_YORK, open: "09:30", close: "16:00"},
                            onCancel: () => {
                                closed.val = true;
                            },
                            onOK: async (newItem) => {
                                const errData = await addSecurity(newItem);
                                if (errData) {
                                    return errData
                                }

                                van.add(sidebar, SidebarItem({item: newItem, state: state}))

                                state.symbol.val = newItem.symbol;
                                state.timeZone.val = newItem.timeZone;
                                state.open.val = newItem.open;
                                state.close.val = newItem.close;
                                state.displayContent.val = true;
                                state.selectedSymbol.val = newItem.symbol;

                                closed.val = true;

                                return null;
                            },
                        }),
                    ),
                ),
            )
        },
    });

    btn.classList.add("fa-solid");
    btn.classList.add("fa-plus");
    btn.classList.add("order-last");
    
    return btn;
}

const SidebarItem = ({item, state}) => {
    const deleted = van.state(false);
    const itemState = {
        symbol: van.state(item.symbol),
        timeZone: van.state(item.timeZone),
        open: van.state(item.open),
        close: van.state(item.close),
    }

    const itemClass = van.derive(() => {
        const isSelected = state.selectedSymbol.val == itemState.symbol.val;

        return `font-semibold md:px-4 md:py-2 ${isSelected ? "text-indigo-500" : "text-gray-500"}`
    });

    return () => deleted.val ? null : button({
        class: itemClass,
        onclick: () => {
            state.selectedSymbol.val = itemState.symbol.val;
            state.editHook.val = () => {
                console.log("editing security %s", itemState.symbol.val);

                const closed = van.state(false);
            
                van.add(
                    document.body,
                    Modal(
                        {closed},
                        div({style: "display: flex; justify-content: center;"},
                            SecurityForm({
                                title: "Edit Security",
                                item: {
                                    symbol: itemState.symbol.val,
                                    timeZone: itemState.timeZone.val,
                                    open: itemState.open.val,
                                    close: itemState.close.val,
                                },
                                onCancel: () => {
                                    closed.val = true;
                                },
                                onOK: async (updatedItem) => {
                                    const errData = await updateSecurity(updatedItem);
                                    if (errData) {
                                        return errData;
                                    }

                                    state.symbol.val = updatedItem.symbol;
                                    state.timeZone.val = updatedItem.timeZone;
                                    state.open.val = updatedItem.open;
                                    state.close.val = updatedItem.close;

                                    itemState.symbol.val = updatedItem.symbol;
                                    itemState.timeZone.val = updatedItem.timeZone;
                                    itemState.open.val = updatedItem.open;
                                    itemState.close.val = updatedItem.close;
                                
                                    closed.val = true;

                                    return null;
                                },
                            }),
                        ),
                    ),
                )
            }
            state.deleteHook.val = async () => {
                if (await delSecurity(item.symbol)) {
                    state.displayContent.val = false;

                    deleted.val = true;
                }
            };
            state.symbol.val = itemState.symbol.val;
            state.timeZone.val = itemState.timeZone.val;
            state.open.val = itemState.open.val;
            state.close.val = itemState.close.val;
            state.displayContent.val = true;
        }},
        item.symbol,
    );
}

const Securities = () => {
    const state = {
        displayContent: van.state(false),
        symbol: van.state(""),
        timeZone: van.state(TZ_NEW_YORK),
        open: van.state("09:30"),
        close: van.state("16:00"),
        editHook: van.state(() => {}),
        deleteHook: van.state(() => {}),
        selectedSymbol: van.state(""),
    };

    const coreHours = van.derive(() => {
        return `${state.open.val} - ${state.close.val}`;
    });
    const contentAreaClass = van.derive(() => {
        return `h-screen flex flex-col px-6 py-4 ${state.displayContent.val ? "" : " hidden"}`;
    });

    const sidebarArea = div(
        {class:"flex flex-col flex-nowrap overflow-y-scroll"},
    )
    const editBtn = ButtonAct({
        text: "",
        onclick: () => state.editHook.val()
    });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => state.deleteHook.val()
    });
    const contentArea = div(
        {class: contentAreaClass},
        div({class: "flex flex-row"}, p({class: "flex grow font-semibold"}, "Symbol"), p(state.symbol)),
        div({class: "flex flex-row"}, p({class: "flex grow font-semibold"}, "Time Zone"), p(state.timeZone)),
        div({class: "flex flex-row"}, p({class: "flex grow font-semibold"}, "Core Hours"), p(coreHours)),
        div({class: "flex flex-row-reverse"}, editBtn, deleteBtn),
    )

    editBtn.classList.add("fa-solid");
    editBtn.classList.add("fa-pen-to-square");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    getSecurities().then(
        (items) => {
            const sidebarItems = items.map(item => SidebarItem({item: item, state: state}));

            van.add(sidebarArea, sidebarItems);
            van.add(sidebarArea, AddNewButton({sidebar: sidebarArea, state: state}));
        }
    );


    return div(
        {class: "h-screen flex"},
        sidebarArea,
        contentArea,
    );
}

export default Securities