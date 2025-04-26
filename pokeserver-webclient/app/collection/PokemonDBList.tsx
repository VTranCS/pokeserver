'use client';

import Image from "next/image";
import {Table, TableHeader, TableColumn, TableBody, TableRow, TableCell,  getKeyValue,} from "@heroui/table";
import { PokemonDBEntry } from "@/app/types/PokemonDBEntry";
export default function PokemonDBList({ pokemonList }: { pokemonList: PokemonDBEntry[] }) {


  console.log("PokemonDBList", pokemonList);
  const columns = [
    {
      key: "Name",
      label: "NAME",
    },
    {
      key: "Id",
      label: "ID",
    },
    {
      key: "Vote",
      label: "VOTE",
    },
    {
      key: "Url",
      label: "SPRITE",
    },
  ];

  const renderCell = (item: PokemonDBEntry, columnKey: string) => {
    if (columnKey === "Url") {
      return <Image src={item.Url} alt={item.Name} width={50} height={50} />;
    }
    else if (columnKey === "Name") {
      return <a
        href={`https://pokemondb.net/pokedex/${item.Id}`}
        className="text-blue-500 hover:underline"
        rel="noopener noreferrer"
        target="_blank"
      >{item.Name}</a>
    }
    else {
      return <p>{item[columnKey as keyof PokemonDBEntry]}</p>
    }
    return getKeyValue(item, columnKey);
  };

  if (!Array.isArray(pokemonList)) {
    console.log(typeof pokemonList["pokemon"]);
    console.error("pokemonList is not an array:");
    return null;
  }
  return (
    <Table isStriped={true} isHeaderSticky aria-label="Example table with dynamic content">
      <TableHeader columns={columns}>
        {(column) => <TableColumn className="text-left" key={column.key}>{column.label}</TableColumn>}
      </TableHeader>
      <TableBody items={pokemonList}>
      {(item) => (
          <TableRow key={item.Id}>
            {(columnKey) => <TableCell>{renderCell(item, columnKey as string)}</TableCell>}
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
