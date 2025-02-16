'use client';

import Image from "next/image";
import React, { useState, useEffect, useRef } from "react";
import { revalidatePath } from 'next/cache';
import Link from "next/link";
import {Table, TableHeader, TableColumn, TableBody, TableRow, TableCell,  getKeyValue,} from "@heroui/table";
export default function Home() {

  interface Pokemon {
    Name: any;
    Id: any;
    Vote: any;
    Url: any;
  }

  let [pokemonList, setPokemonList] = useState<Pokemon[]>([]);
  const [fetchTrigger, setFetchTrigger] = useState(true); // State to trigger useEffect

  useEffect(() => {

    const fetchRandomPokemon = async () => {
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

    fetchRandomPokemon();
  }, [fetchTrigger]); // Only runs when fetchTrigger changes

  const [textInput, setTextInput] = useState('');

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
      label: "URL",
    },
  ];
  return (
    <Table aria-label="Example table with dynamic content">
      <TableHeader columns={columns}>
        {(column) => <TableColumn key={column.key}>{column.label}</TableColumn>}
      </TableHeader>
      <TableBody items={pokemonList}>
        {(item) => (
          <TableRow key={item.Id}>
            {(columnKey) => <TableCell>{getKeyValue(item, columnKey)}</TableCell>}
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}
