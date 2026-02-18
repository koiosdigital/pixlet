import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import DeleteIcon from '@mui/icons-material/Delete';

import { MapContainer, TileLayer, Polygon, Polyline, CircleMarker, Marker, useMapEvents } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

import { set } from '../../../config/configSlice';
import styles from './styles.css';

// Fix Leaflet default marker icon paths for webpack
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png').default,
    iconUrl: require('leaflet/dist/images/marker-icon.png').default,
    shadowUrl: require('leaflet/dist/images/marker-shadow.png').default,
});

// Convert internal [lat, lng] positions to GeoJSON [lng, lat] Polygon
function toGeoJSONPolygon(positions) {
    if (positions.length < 3) return null;
    const coords = positions.map(p => [p[1], p[0]]);
    coords.push(coords[0]); // close the ring
    return {
        type: 'Polygon',
        coordinates: [coords],
    };
}

// Convert GeoJSON [lng, lat] coordinates to internal [lat, lng] positions
function fromGeoJSONPolygon(geojson) {
    if (!geojson || !geojson.coordinates || !geojson.coordinates[0]) return [];
    const ring = geojson.coordinates[0];
    const coords = ring.slice(0, ring.length - 1);
    return coords.map(c => [c[1], c[0]]);
}

// Build a GeoJSON FeatureCollection from the current state
function buildValue(point, positions, closed, collectPoint) {
    const features = [];
    if (collectPoint && point) {
        features.push({
            type: 'Feature',
            properties: { role: 'point' },
            geometry: { type: 'Point', coordinates: [point[1], point[0]] },
        });
    }
    if (closed && positions.length >= 3) {
        features.push({
            type: 'Feature',
            properties: { role: 'polygon' },
            geometry: toGeoJSONPolygon(positions),
        });
    }
    if (features.length === 0) return '';
    return JSON.stringify({ type: 'FeatureCollection', features });
}

// Parse initial state from a stored FeatureCollection config value
function parseInitialState(valueStr, collectPoint) {
    if (!valueStr) return { point: null, positions: [] };
    try {
        const parsed = JSON.parse(valueStr);
        if (parsed.type === 'FeatureCollection' && parsed.features) {
            let point = null;
            let positions = [];
            for (const feature of parsed.features) {
                if (!feature.geometry) continue;
                if (feature.geometry.type === 'Point') {
                    const c = feature.geometry.coordinates;
                    point = [c[1], c[0]]; // [lat, lng]
                } else if (feature.geometry.type === 'Polygon') {
                    positions = fromGeoJSONPolygon(feature.geometry);
                }
            }
            return { point, positions };
        }
    } catch {}
    return { point: null, positions: [] };
}

function ClickHandler({ onClick }) {
    useMapEvents({
        click(e) {
            onClick(e.latlng);
        },
    });
    return null;
}

export default function GeoJSONField({ field }) {
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();
    const collectPoint = field.collect_point || false;

    const getInitialValue = () => {
        if (field.id in config && config[field.id].value) {
            return config[field.id].value;
        }
        return field.default || '';
    };

    const initial = parseInitialState(getInitialValue(), collectPoint);

    const [point, setPoint] = useState(initial.point);
    const [positions, setPositions] = useState(initial.positions);
    const [closed, setClosed] = useState(initial.positions.length >= 3);
    // Step: 0 = select point (location), 1 = draw polygon
    const [step, setStep] = useState(() => {
        if (!collectPoint) return 1;
        return initial.point ? 1 : 0;
    });

    // Sync to Redux config when state changes
    useEffect(() => {
        const value = buildValue(point, positions, closed, collectPoint);
        if (value) {
            dispatch(set({ id: field.id, value }));
        }
    }, [point, positions, closed]);

    // Set default on mount if no config value exists
    useEffect(() => {
        if (!(field.id in config) && field.default) {
            dispatch(set({ id: field.id, value: field.default }));
        }
    }, []);

    const handleMapClick = (latlng) => {
        if (collectPoint && step === 0) {
            // Step 1: place the user's location point
            setPoint([latlng.lat, latlng.lng]);
            setStep(1);
        } else if (step === 1 && !closed) {
            // Step 2: add polygon vertices
            setPositions([...positions, [latlng.lat, latlng.lng]]);
        }
    };

    const handleClosePolygon = () => {
        if (positions.length >= 3) {
            setClosed(true);
        }
    };

    const handleClear = () => {
        setPoint(null);
        setPositions([]);
        setClosed(false);
        setStep(collectPoint ? 0 : 1);
        dispatch(set({ id: field.id, value: '' }));
    };

    const handleClearPolygon = () => {
        setPositions([]);
        setClosed(false);
    };

    const getCenter = () => {
        if (positions.length > 0) {
            const avgLat = positions.reduce((s, p) => s + p[0], 0) / positions.length;
            const avgLng = positions.reduce((s, p) => s + p[1], 0) / positions.length;
            return [avgLat, avgLng];
        }
        if (point) return point;
        return [40.678, -73.944];
    };

    const getStatusText = () => {
        if (collectPoint && step === 0) {
            return 'Click the map to set your location';
        }
        if (closed) {
            return `Polygon complete (${positions.length} points)`;
        }
        if (positions.length === 0) {
            return 'Click on the map to start drawing a polygon';
        }
        if (positions.length < 3) {
            return `Click to add more points (${positions.length}/3 minimum)`;
        }
        return 'Click to add points, or press "Close Polygon" to finish';
    };

    return (
        <div>
            {collectPoint && (
                <Stepper activeStep={step} style={{ marginBottom: '12px' }}>
                    <Step completed={point !== null}>
                        <StepLabel>Set Location</StepLabel>
                    </Step>
                    <Step completed={closed}>
                        <StepLabel>Draw Polygon</StepLabel>
                    </Step>
                </Stepper>
            )}
            <div className={styles.mapContainer}>
                <MapContainer
                    center={getCenter()}
                    zoom={positions.length > 0 || point ? 13 : 4}
                    style={{ height: '300px', width: '100%' }}
                >
                    <TileLayer
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    <ClickHandler onClick={handleMapClick} />
                    {point && (
                        <Marker position={point} />
                    )}
                    {positions.length >= 3 && closed && (
                        <Polygon
                            positions={positions}
                            pathOptions={{ color: '#1976d2', fillOpacity: 0.2 }}
                        />
                    )}
                    {positions.length >= 2 && !closed && (
                        <Polyline
                            positions={positions}
                            pathOptions={{ color: '#1976d2' }}
                        />
                    )}
                    {positions.map((pos, i) => (
                        <CircleMarker
                            key={i}
                            center={pos}
                            radius={5}
                            pathOptions={{
                                color: '#1976d2',
                                fillColor: i === 0 && positions.length >= 3 && !closed ? '#4caf50' : '#1976d2',
                                fillOpacity: 1,
                            }}
                        />
                    ))}
                </MapContainer>
            </div>
            <Typography variant="body2" style={{ margin: '8px 0' }}>
                {getStatusText()}
            </Typography>
            <Stack spacing={2} direction="row">
                {!closed && positions.length >= 3 && (
                    <Button
                        variant="contained"
                        onClick={handleClosePolygon}
                        size="small"
                    >
                        Close Polygon
                    </Button>
                )}
                {collectPoint && positions.length > 0 && !closed && (
                    <Button
                        variant="outlined"
                        onClick={handleClearPolygon}
                        size="small"
                    >
                        Restart Polygon
                    </Button>
                )}
                {(positions.length > 0 || point) && (
                    <Button
                        variant="outlined"
                        onClick={handleClear}
                        startIcon={<DeleteIcon />}
                        size="small"
                    >
                        Clear All
                    </Button>
                )}
            </Stack>
        </div>
    );
}
