'use client';

import Image from "next/image";
import {Table, TableHeader, TableColumn, TableBody, TableRow, TableCell} from "@heroui/table";
import { FetchedPokemonDBEntry } from "@/app/types/PokemonDBEntry";
export default function PokemonDBList({ pokemonList }: { pokemonList: FetchedPokemonDBEntry[] }) {


  console.log("PokemonDBList", pokemonList);
  const columns = [
    {
      key: "name",
      label: "POKEMON",
    },
    {
      key: "id",
      label: "ID",
    },
    {
      key: "vote",
      label: "VOTE",
    },
    {
      key: "Url",
      label: "",
    },
  ];

  const renderCell = (item: FetchedPokemonDBEntry, columnKey: string) => {
    if (columnKey === "Url") {
      return ;
    }
    if (columnKey === "name") {
      return <a
        href={`https://pokemondb.net/pokedex/${item.id}`}
        className="text-blue-500 hover:underline"
        rel="noopener noreferrer"
        target="_blank"
      ><Image src={item.image} alt={item.name} width={100} height={100} />{item.name}</a>
    }
    else {
      return <p>{item[columnKey as keyof FetchedPokemonDBEntry]}</p>
    }
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
          <TableRow key={item.id}>
            {(columnKey) => <TableCell>{renderCell(item, columnKey as string)}</TableCell>}
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
