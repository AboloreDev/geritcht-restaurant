import React from "react";
import Navbar from "./homepage/Navbar";

interface Children {
  children: React.ReactNode;
}

const HomepageLayout = ({ children }: Children) => {
  return (
    <div>
      <Navbar />
      <main>{children}</main>
    </div>
  );
};

export default HomepageLayout;
