import React from "react";
import MenuDetailContent from "./components/MenuDetailContent";
import EditMenuSheet from "../components/edit/EditSheet";
import DeleteMenuDialog from "../components/delete/DeleteMenuDialog";

const MenuDetailsPage = () => {
  return (
    <div>
      <MenuDetailContent />

      <EditMenuSheet />

      <DeleteMenuDialog />
    </div>
  );
};

export default MenuDetailsPage;
