type User = {
  name: string;
  age: number;
  email: string;
};

const users: User[] = [
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
];

const users2: { name: string; age: number; email: string }[] = [
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
  { name: 'lorem', age: 0, email: 'lorem ipsum' },
];
