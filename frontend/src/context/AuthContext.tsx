"use client";

import { createContext, useContext, useEffect, useState } from "react";
import Keycloak from "keycloak-js";

const AuthContext = createContext<{
  authenticated: boolean;
  token: string | null;
  user: any;
} | null>(null);

const keycloakConfig = {
  url: "http://localhost:8080",
  realm: "omniguard",
  clientId: "omniguard-frontend",
};

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false);
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    const kc = new Keycloak(keycloakConfig);
    kc.init({ onLoad: "login-required", checkLoginIframe: false }).then((auth) => {
      setAuthenticated(auth);
      setToken(kc.token || null);
      setUser(kc.tokenParsed);

      // Keep token fresh
      setInterval(() => {
        kc.updateToken(70).then((refreshed) => {
          if (refreshed) {
            setToken(kc.token || null);
          }
        });
      }, 60000);
    }).catch(() => {
      console.error("Authenticated Failed");
    });
  }, []);

  return (
    <AuthContext.Provider value={{ authenticated, token, user }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
