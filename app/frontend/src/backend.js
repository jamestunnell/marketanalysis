const BASE_URL = '/api';

const Delete = async (route) => {
    const url = BASE_URL + route;

    console.log("DELETE %s", url);

    return fetch(url, {method: "DELETE"});
}

const Get = async (route) => {
    const url = BASE_URL + route;

    console.log("GET %s", url);

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

    console.log("POST %s", url);

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

const PutJSON = async ({route, object}) => {
    const url = BASE_URL + route;

    console.log("PUT %s", url);

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

export {Delete, Get, PostJSON, PutJSON};