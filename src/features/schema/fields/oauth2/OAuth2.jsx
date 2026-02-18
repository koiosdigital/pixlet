import React, { useState, useEffect, useRef } from "react";
import { useDispatch, useSelector } from "react-redux";
import OAuth2Login from "react-simple-oauth2-login";

import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";

import { callHandlerSetValue } from "../../../handlers/actions";
import { set as setError } from "../../../errors/errorSlice";
import { set, remove } from "../../../config/configSlice";

function generateCodeVerifier() {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}

async function generateCodeChallenge(verifier) {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const digest = await crypto.subtle.digest("SHA-256", data);
  return base64UrlEncode(new Uint8Array(digest));
}

function base64UrlEncode(bytes) {
  return btoa(String.fromCharCode(...bytes))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

export default function OAuth2({ field }) {
  const [loggedIn, setLoggedIn] = useState("");
  const [pkceParams, setPkceParams] = useState(null);
  const [userClientId, setUserClientId] = useState("");
  const [userClientSecret, setUserClientSecret] = useState("");
  const codeVerifierRef = useRef(null);
  const dispatch = useDispatch();
  const config = useSelector((state) => state.config);
  const host = document.location.host.replaceAll("localhost", "127.0.0.1");
  const redirectUri =
    document.location.protocol + "//" + host + "/oauth-callback";
  const clientId = field.user_defined_client ? userClientId : field.client_id;

  useEffect(() => {
    if (field.id in config) {
      setLoggedIn(config[field.id].value);
    }
  }, [config]);

  useEffect(() => {
    if (field.pkce) {
      const verifier = generateCodeVerifier();
      codeVerifierRef.current = verifier;
      generateCodeChallenge(verifier).then((challenge) => {
        setPkceParams({
          code_challenge_method: "S256",
          code_challenge: challenge,
        });
      });
    }
  }, [field.pkce]);

  const onSuccess = (response) => {
    if (!response.code) {
      return onFailure("access was not granted");
    }

    const handlerParams = {
      code: response.code,
      client_id: clientId,
      redirect_uri: redirectUri,
      grant_type: "authorization_code",
    };

    if (field.pkce && codeVerifierRef.current) {
      handlerParams.code_verifier = codeVerifierRef.current;
    }

    if (field.user_defined_client && userClientSecret) {
      handlerParams.client_secret = userClientSecret;
    }

    callHandlerSetValue(
      field.id,
      field.handler,
      handlerParams,
      (value) => {
        setLoggedIn(value);
        dispatch(
          set({
            id: field.id,
            value: value,
          }),
        );

        // Regenerate PKCE params for next login attempt
        if (field.pkce) {
          const verifier = generateCodeVerifier();
          codeVerifierRef.current = verifier;
          generateCodeChallenge(verifier).then((challenge) => {
            setPkceParams({
              code_challenge_method: "S256",
              code_challenge: challenge,
            });
          });
        }
      },
    );
  };

  const logout = () => {
    setLoggedIn("");
    dispatch(remove(field.id));
  };

  const onFailure = (response) => {
    let msg = `failed login: ${response}`;
    dispatch(setError({ id: msg, message: msg }));
    console.error(response);
  };

  const renderButton = (params) => {
    return (
      <Button variant="contained" onClick={params.onClick}>
        Login
      </Button>
    );
  };

  if (loggedIn) {
    return (
      <Button variant="contained" onClick={logout}>
        Logout
      </Button>
    );
  }

  // Wait for PKCE params to be generated before rendering
  if (field.pkce && !pkceParams) {
    return (
      <Button variant="contained" disabled>
        Login
      </Button>
    );
  }

  let scope = field.scopes.join(" ");

  return (
    <Box sx={{ display: "flex", gap: 1, alignItems: "center" }}>
      {field.user_defined_client && (
        <>
          <TextField
            size="small"
            label="Client ID"
            value={userClientId}
            onChange={(e) => setUserClientId(e.target.value)}
          />
          <TextField
            size="small"
            label="Client Secret"
            type="password"
            value={userClientSecret}
            onChange={(e) => setUserClientSecret(e.target.value)}
          />
        </>
      )}
      {clientId ? (
        <OAuth2Login
          isCrossOrigin={true}
          authorizationUrl={field.authorization_endpoint}
          responseType="code"
          scope={scope}
          state="abc123"
          clientId={clientId}
          redirectUri={redirectUri}
          render={renderButton}
          onSuccess={onSuccess}
          onFailure={onFailure}
          {...(pkceParams && { extraParams: pkceParams })}
        />
      ) : (
        <Button variant="contained" disabled>
          Login
        </Button>
      )}
    </Box>
  );
}
