import { useContext, createContext } from "react";
import {
  setPersistence,
  inMemoryPersistence,
  signInWithPopup,
  signOut,
} from "firebase/auth";
import { auth } from "@/firebase/clientApp";
import { useRouter } from "next/router";
import { useLocalStorage } from "./useLocalStorage";

const AuthContext = createContext();

const createSession = async (idToken) => {
  const axios = require("axios");
  const clientId = Buffer.from(process.env.NEXT_PUBLIC_CLIENT_ID).toString(
    "base64"
  );
  return axios.post(process.env.NEXT_PUBLIC_APP_HOST + "/api/auth/sign-in", {
    client_id: clientId,
    id_token: idToken,
  });
};

const endSession = async () => {
  const axios = require("axios");
  return axios.post(process.env.NEXT_PUBLIC_APP_HOST + "/api/auth/sign-out", {
    withCredentials: true,
  });
};

export const AuthContextProvider = ({ children }) => {
  const [user, setUser] = useLocalStorage("user");
  const router = useRouter();

  const firebaseSignIn = (provider) => {
    setPersistence(auth, inMemoryPersistence);
    signInWithPopup(auth, provider)
      .then((result) => {
        return result.user.getIdToken().then((idToken) => {
          createSession(idToken)
            .then(() => {
              setUser({
                ID: result.user.uid,
                displayName: result.user.displayName,
                photoURL: result.user.photoURL,
              });
              router.push("/messages");
            })
            .catch((error) => console.error(error));
        });
      })
      .catch((error) => console.error(error.code, error.message));
  };

  const firebaseSignOut = () => {
    signOut(auth).then(() => {
      endSession().then(() => {
        setUser(null);
        router.push(process.env.NEXT_PUBLIC_APP_HOST);
      });
    });
  };

  return (
    <AuthContext.Provider value={{ firebaseSignIn, firebaseSignOut, user }}>
      {children}
    </AuthContext.Provider>
  );
};

export const UserAuth = () => {
  return useContext(AuthContext);
};
