"use client";

import { createContext, useContext } from "react";
import {
  setPersistence,
  signInWithPopup,
  signOut,
  inMemoryPersistence,
} from "firebase/auth";
import { auth } from "@/lib/firebase";
import { useRouter } from "next/navigation";
import { useLocalStorage } from "@/hooks/useLocalStorage";
import { CreateSession, EndSession } from "@/lib/authHelper";

const AuthContext = createContext();
export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const router = useRouter();
  const [user, setUser] = useLocalStorage("user");

  const firebaseSignIn = (provider) => {
    setPersistence(auth, inMemoryPersistence);
    signInWithPopup(auth, provider)
      .then((result) => {
        return result.user.getIdToken().then((idToken) => {
          CreateSession(idToken)
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
      EndSession().then(() => {
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

export default AuthProvider;
