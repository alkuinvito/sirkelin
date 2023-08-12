import firebase from "firebase/compat/app";
import { getAuth } from "firebase/auth";

const app = firebase.initializeApp({
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
});

export const auth = getAuth(app);
