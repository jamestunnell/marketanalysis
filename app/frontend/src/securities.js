import van from "vanjs-core"
import { routeTo } from 'vanjs-router'
import {Modal} from "vanjs-ui"

import {Delete, Get, PostJSON, PutJSON} from './backend.js'
import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from "./buttons.js";
import { ButtonGroup } from './buttongroup.js'
import { IconAdd, IconClose, IconDelete, IconEdit, IconRefresh, IconView } from './icons.js'
import { INPUT_CLASS } from './input.js'
import {Table, TableRow} from './table.js'
import { ModalBackground } from "./modal.js";
import { AppErrorAlert } from "./apperror.js";

const {div, input, label, option, p, select, tbody} = van.tags

const TIME_ZONES = Intl.supportedValuesOf('timeZone');
const TZ_NEW_YORK = "America/New_York"
const MAX_DAYS = 100

const getSecurities = async () => {
    console.log("getting securities");

    const resp = await Get('/securities');

    if (resp.status != 200) {
        console.log("failed to get securities", await resp.json());

        return []
    }

    const d = await resp.json();

    console.log(`received ${d.securities.length} securities`, d.securities);

    return d.securities;
}

const createSecurity = async (item) => {
    return new Promise((resolve, reject) => {
        console.log("creating security", item);

        PostJSON({route: '/securities', object: item}).then(resp => {
            if (resp.status != 204) {
                resp.json().then(appErr => {
                    console.log("failed to create security", appErr);

                    reject(appErr)
                })
        
                return;
            }

            console.log(`created security %s`, item.symbol);

            // Avoid Fetch failed loading
            resp.text().then(text => {
                resolve()
            })
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: "failed to create security",
                details: [err.message],
            })
        })
    })
}

const updateSecurity = async (item) => {
    return new Promise((resolve, reject) => {
        console.log("updating security", item);

        PutJSON({route: `/securities/${item.symbol}`, object: item}).then(resp => {
            if (resp.status != 204) {
                resp.json().then(appErr => {
                    console.log("failed to update security", appErr);

                    reject(appErr)
                }).catch((err) => {
                    reject({
                        title: "Action Failed",
                        message: "failed to update security",
                        details: [resp.statusText, err.message],
                    })
                })
        
                return;
            }

            console.log(`updated security ${item.symbol}`);

            // Avoid Fetch failed loading
            resp.text().then(text => {
                resolve()
            })
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: "failed to update security",
                details: [err.message],
            })
        })
    })
}

const deleteSecurity = async (symbol) => {
    console.log("deleting security", symbol);

    const resp = await Delete(`/securities/${symbol}`);

    if (resp.status != 204) {
        console.log("failed to delete security", await resp.json());

        return false
    }

    // Avosymbol Fetch failed loading
    await resp.text();

    console.log(`deleted security %s`, symbol);

    return true;
}

const SecurityModal = ({obj, actionName, onOK}) => {
    const closed = van.state(false)
    const symbol = van.state(obj.symbol)
    const timeZone = van.state(obj.timeZone)
    const days = van.state(obj.days)
    const incomplete = van.derive(() => {
        if (symbol.val.length === 0) {
            return true
        }

        if (timeZone.val.length === 0) {
            return true
        }

        if (days.val < 1) {
            return true
        }
    })

    const closeBtn = ButtonIcon({
        icon: IconClose(),
        onclick: ()=> closed.val = true,
    })
    
    closeBtn.classList.add("self-end")

    const timeZoneOpts = TIME_ZONES.map(x => {
        let props = {
            value: x,
            selected: (x === obj.timeZone),
        }

        return option(props, x)
    })
    const modal = ModalBackground(
        div(
            {id: "foreground", class: "flex flex-col space-y-3 block p-16 rounded-lg bg-white min-w-[25%] max-w-[50%]"},
            closeBtn,
            p({class: "text-lg font-medium font-bold text-center"}, `${actionName} Security`),
            label({for: "symbol"}, "Symbol"),
            input({
                id: "symbol",
                class: INPUT_CLASS,
                type: "text",
                value: obj.symbol,
                placeholder: "Unique, non-empty (e.g. SPY, QQQ)",
                required: true,
                disabled: obj.symbol.length > 0,
                oninput: e => symbol.val = e.target.value,
            }),
            label({for: "timeZone"}, "Time Zone"),
            select({
                id: "timeZone",
                class: INPUT_CLASS,
                onchange: e => timeZone.val = e.target.value,
            }, timeZoneOpts),
            label({for:"days"}, "Days to Collect"),
            input({
                id: "days",
                type: "number",
                class: INPUT_CLASS,
                value: days, 
                min: 1,
                step: 1,
                onchange: e => {
                    const newVal = parseInt(e.target.value, 10)
                    if (isNaN(newVal)) {
                        days.val = 0    

                        return false
                    }
                    
                    days.val = newVal
                },
            }),
            div(
                {class:"mt-4 flex justify-center"},
                ButtonCancel({
                    child: "Cancel",
                    onclick: ()=> closed.val = true,
                }),
                Button({
                    child: "OK",
                    disabled: incomplete,
                    onclick: ()=> {
                        onOK({
                            symbol: symbol.val,
                            timeZone: timeZone.val,
                            days: days.val,
                        })

                        closed.val = true
                    },
                }),
            ),
        )
    )

    van.add(document.body, () => closed.val ? null : modal);
}

const SecurityRow = (obj) => {
    const deleted = van.state(false)
    const symbol = van.state(obj.symbol)
    const timeZone = van.state(obj.timeZone)
    const days = van.state(obj.days)

    // const viewBtn = ButtonIcon({
    //     icon: IconView(),
    //     // text: "View",
    //     onclick: () => routeTo('securities', [symbol.val]),
    // });
    const editBtn = ButtonIcon({
        icon: IconEdit(),
        // text: "Edit",
        onclick: () => {
            SecurityModal({
                obj: {symbol: symbol.val, timeZone: timeZone.val, days: days.val},
                actionName: "Edit",
                onOK: (obj) => {
                    updateSecurity(obj).then(() => {
                      symbol.val = obj.symbol
                      timeZone.val = obj.timeZone
                      days.val = obj.days
                    }).catch(appErr => AppErrorAlert(appErr))
                },
            })
        },
    });
    const deleteBtn = ButtonIcon({
        icon: IconDelete(),
        // text: "Delete",
        onclick: () => {
            deleteSecurity(symbol.val).then(ok => {
                if (ok) {
                    deleted.val = true
                }
            })
        },
    });

    const actionButtons = ButtonGroup({buttons: [editBtn, deleteBtn]});
    const rowItems = [symbol, timeZone, days, actionButtons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const SecuritiesPage = () => {
    const columnNames = ["Symbol", "Time Zone", "Days", ""]
    const tableBody = tbody({class:"table-auto"});

    getSecurities().then((objs) => {
        const rows = objs.map(obj => SecurityRow(obj));

        van.add(tableBody, rows);
    });

    const addIcon = IconAdd()
    
    addIcon.classList.add("text-xl")

    const newSecurityBtn = ButtonIcon({
        icon: addIcon,
        onclick: () => {
            SecurityModal({
                obj: {symbol: "", timeZone: TZ_NEW_YORK, days: 20},
                actionName: "Add",
                onOK: (obj)=> {
                    createSecurity(obj).then(() => {
                        van.add(tableBody, SecurityRow(obj));
                    }).catch(appErr => AppErrorAlert(appErr));
                },
            })
        },
    });

    return div(
        {class: "container p-4 w-full flex flex-col divide-y divide-gray-400"},
        div(
            {class: "flex flex-col mt-4"},
            div(
                {class: "grid grid-cols-2"},
                div(
                    {class: "flex flex-row p-2"},
                    p({class: "p-3 m-1 text-xl font-medium"}, "Securities"),
                ),
                div(
                    {class: "flex flex-row-reverse p-2"},
                    newSecurityBtn,
                )
            ),
            Table({columnNames: columnNames, tableBody: tableBody}),
        )
    )
}

export default SecuritiesPage;