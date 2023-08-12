import "./globals.css";
import { Poppins } from "next/font/google";
import StyledComponentsRegistry from "@/lib/AntdRegistry";
import AuthProvider from "@/contexts/AuthContext";
import Providers from "@/contexts/QueryProvider";

const poppins = Poppins({ subsets: ["latin"], weight: "400" });

export const metadata = {
  title: "Sirkelin",
  description: "Centralized messaging platform",
};

export default function RootLayout({ children }) {
  return (
    <AuthProvider>
      <html lang="en">
        <body className={poppins.className}>
          <Providers>
            <StyledComponentsRegistry>{children}</StyledComponentsRegistry>
          </Providers>
        </body>
      </html>
    </AuthProvider>
  );
}
