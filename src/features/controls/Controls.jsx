import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import axios from 'axios';

import { Button, Stack, TextField } from '@mui/material';
import { resetConfig, setConfig } from '../config/actions';
import { set } from '../config/configSlice';

export default function Controls() {
    const preview = useSelector(state => state.preview);
    const config = useSelector(state => state.config);
    const schema = useSelector(state => state.schema);
    const dispatch = useDispatch();

    const [width, setWidth] = useState(64);
    const [height, setHeight] = useState(32);

    useEffect(() => {
        axios.get(`${PIXLET_API_BASE}/api/v1/dimensions`)
            .then(res => {
                setWidth(res.data.width);
                setHeight(res.data.height);
            })
            .catch(() => {});
    }, []);

    function setDimensions() {
        axios.post(`${PIXLET_API_BASE}/api/v1/dimensions`, { width, height })
            .catch(() => {});
    }

    let imageType = 'webp';
    if (preview.value.img_type === "gif") {
        imageType = 'gif';
    }

    function downloadPreview() {
        const date = new Date().getTime();
        const element = document.createElement("a");

        // convert base64 to raw binary data held in a string
        let byteCharacters = atob(preview.value.img);

        // create an ArrayBuffer with a size in bytes
        let arrayBuffer = new ArrayBuffer(byteCharacters.length);

        // create a new Uint8Array view
        let uint8Array = new Uint8Array(arrayBuffer);

        // assign the values
        for (let i = 0; i < byteCharacters.length; i++) {
            uint8Array[i] = byteCharacters.charCodeAt(i);
        }

        const file = new Blob([uint8Array], { type: 'image/' + imageType });
        element.href = URL.createObjectURL(file);
        element.download = `tidbyt-preview-${date}.${imageType}`;
        document.body.appendChild(element); // Required for this to work in FireFox
        element.click();
    }

    function downloadConfig() {
        const date = new Date().getTime();
        const element = document.createElement("a");
        const jsonData = config;

        // Use Blob object for JSON
        const file = new Blob([JSON.stringify(jsonData)], { type: 'application/json' });
        element.href = URL.createObjectURL(file);
        element.download = `config-${date}.json`;
        document.body.appendChild(element); // Required for this to work in FireFox
        element.click();
    }

    function selectConfig() {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'application/json';

        input.onchange = function (event) {
            const file = event.target.files[0];
            if (file.type !== "application/json") {
                return;
            }

            const reader = new FileReader();

            reader.onload = function () {
                let contents = reader.result;
                let json = JSON.parse(contents);
                setConfig(json);
            };

            reader.onerror = function () {
                console.log(reader.error);
            };

            reader.readAsText(file);
        };

        input.click();
    }


    function resetSchema() {
        history.replaceState(null, '', location.pathname);
        resetConfig();
        schema.value.schema.forEach((field) => {
            if (field.default) {
                dispatch(set({
                    id: field.id,
                    value: field.default,
                }));
            };
        });
    };

    return (
        <Stack sx={{ marginTop: '32px' }} spacing={2} direction="row" alignItems="center">
            <Button variant="outlined" onClick={() => selectConfig()}>Open Config</Button>
            <Button variant="outlined" onClick={() => downloadConfig()}>Save Config</Button>
            <Button variant="outlined" onClick={() => resetSchema()}>Reset</Button>
            <Button variant="contained" onClick={() => downloadPreview()}>Export Image</Button>
            <TextField
                label="Width"
                type="number"
                size="small"
                value={width}
                onChange={(e) => setWidth(parseInt(e.target.value) || 0)}
                sx={{ width: 80 }}
                inputProps={{ min: 1 }}
            />
            <TextField
                label="Height"
                type="number"
                size="small"
                value={height}
                onChange={(e) => setHeight(parseInt(e.target.value) || 0)}
                sx={{ width: 80 }}
                inputProps={{ min: 1 }}
            />
            <Button variant="contained" onClick={() => setDimensions()}>Set Dimensions</Button>
        </Stack>
    );
}