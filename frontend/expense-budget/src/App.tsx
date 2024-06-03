"use client"
import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
// import { useForm } from 'react-hook-form'
// import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from './components/ui/form'
// import { Input } from './components/ui/input'
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import axios from "axios";
const formSchema = z
  .object({
    email: z.string().email(),
    password: z.string().min(3),
  
  })

function App() {

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });



  const handleSubmit = (values: z.infer<typeof formSchema>) => {
    axios.post("http://localhost:8080/login", values,{withCredentials: true}).then((res) => {
      localStorage.setItem("token", res.data.token);
      localStorage.setItem("user",res.data.data.first_name)
      console.log(res);
    });
  };

  return (
    <>
      
   
      <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(handleSubmit)}
          className="max-w-md w-full flex flex-col gap-4"
        >
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => {
              return (
                <FormItem>
                  <FormLabel>Email address</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Email address"
                      type="email"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              );
            }}
          />
          
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => {
              return (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input placeholder="Password" type="password" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              );
            }}
          />
          <Button type="submit" className="w-full">
            Submit
          </Button>
        </form>
      </Form>
    </main>

   </>
  )
}

export default App
