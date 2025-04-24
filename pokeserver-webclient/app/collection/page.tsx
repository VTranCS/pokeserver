'use client';

import Image from "next/image";
import React, { useState, useEffect} from "react";
import {Table, TableHeader, TableColumn, TableBody, TableRow, TableCell,  getKeyValue,} from "@heroui/table";
export default function Home() {

  interface Pokemon {
    Name: string;
    Id: number;
    Vote: number;
    Url: string;
  }

  const [pokemonList, setPokemonList] = useState<Pokemon[]>([]);
  const [fetchTrigger, setFetchTrigger] = useState(true); // State to trigger useEffect

  useEffect(() => {

    const fetchAllPokemon = async () => {
      console.log('url', process.env);
      fetch((process.env.NEXT_PUBLIC_POKEMON_BASE_URL || 'localhost:3000')
        + (process.env.NEXT_PUBLIC_SUFFIX_GET_POKEMON_VOTES || '/getall'))
        .then(response => response.json())
        .then(data => {
          console.log(data);
          setPokemonList(data.pokemon)
        })
      setFetchTrigger(false); // Reset trigger
    };

    fetchAllPokemon();
  }, [fetchTrigger]); 

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

  const renderCell = (item: Pokemon, columnKey: any) => {
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
      return <p>{item[columnKey as keyof Pokemon]}</p>
    }
    return getKeyValue(item, columnKey);
  };
  return (
    <Table isStriped={true} isHeaderSticky aria-label="Example table with dynamic content">
      <TableHeader columns={columns}>
        {(column) => <TableColumn className="text-left" key={column.key}>{column.label}</TableColumn>}
      </TableHeader>
      <TableBody items={pokemonList}>
      {(item) => (
          <TableRow key={item.Id}>
            {(columnKey) => <TableCell>{renderCell(item, columnKey)}</TableCell>}
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
