import React, { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import { callGeneratedHandler } from '../../handlers/actions';
import { updateGenerated } from '../schemaSlice';
import { set as setError } from '../../errors/errorSlice';


export default function Generated({ field }) {
    const [source, setSource] = useState(null);
    const config = useSelector(state => state.config);
    const schema = useSelector(state => state.schema);
    const dispatch = useDispatch();

    useEffect(() => {
        setSource(getSourceField());
    }, [schema])

    useEffect(() => {
        if (!source) return;
        if (source.id in config) {
            callGeneratedHandler(field.id, field.handler, config[source.id].value);
        } else if (source.default !== undefined && source.default !== '') {
            callGeneratedHandler(field.id, field.handler, source.default);
        } else {
            // Source field was removed (e.g. OAuth2 logout) with no default.
            // Clear the generated fields so stale UI doesn't linger.
            dispatch(updateGenerated({ version: '1', schema: [] }));
        }
    }, [config, source])

    const getSourceField = () => {
        if (schema.value.schema.length == 0) {
            return null;
        }

        for (let i = 0; i < schema.value.schema.length; i++) {
            if (schema.value.schema[i].id === field.source) {
                return schema.value.schema[i];
            }
        }

        let msg = `schema.Generated references source that does not exist: ${field.source}`;
        dispatch(setError({ id: msg, message: msg }));
        return null;
    }

    return null;
}