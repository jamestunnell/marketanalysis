const BASE_URL = '/api';

const Delete = async (route) => {
    const url = BASE_URL + route;

    console.log("fetch DELETE %s", url);

    return fetch(url, {method: "DELETE"});
}

const Get = async (route) => {
    const url = BASE_URL + route;

    console.log("fetch GET %s", url);

    return fetch(
        url,
        {
            method: "GET",
            headers: {
                'Accept': 'application/json',
            }
        },
    );
}

const PostJSON = async ({route, object, options={}}) => {
    const url = BASE_URL + route;
    const body = JSON.stringify(object);

    console.log("fetch POST %s with JSON %s", url, body);

    return fetch(
        url,
        {
            method: "POST",
            headers: {
                'Content-Type': 'application/json; charset=utf-8'
            },
            body: body,
            ...options,
        },
    );
}

const Put = async ({route, object}) => {
    const url = BASE_URL + route;

    console.log("fetch PUT %s", url);

    return fetch(
        url,
        {
            method: "PUT",
            headers: {
                'Content-Type': 'application/json; charset=utf-8'
            },
            body: JSON.stringify(object),
        },
    );
}

export {Delete, Get, PostJSON, Put};