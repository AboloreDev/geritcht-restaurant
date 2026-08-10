import React from "react";
import TableDetailContent from "./components/TableDetails";
import EditTableSheet from "../components/EditTableSheet";
import DeleteTableDialog from "../components/DeleteTableDialog";

const page = () => {
  return (
    <div>
      <TableDetailContent />
      <EditTableSheet />
      <DeleteTableDialog />
    </div>
  );
};

export default page;
