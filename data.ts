type User = {
    name: string;
    surname: string;
    age: number;
    email: string;
    premiumAccount: boolean;
    role: string;
}

const users: User[] = [
    {
        name: `John`,
        surname: `Johnson`,
        age: 93,
        email: `sushisamurai@email.com`,
        premiumAccount: false,
        role: `user`,
    },
    {
        name: `Richard`,
        surname: `Williams`,
        age: 91,
        email: `cosmicjellybean@email.com`,
        premiumAccount: true,
        role: `mod`,
    },
    {
        name: `Karen`,
        surname: `Harris`,
        age: 28,
        email: `techno-unicorn@email.com`,
        premiumAccount: true,
        role: `user`,
    },
    {
        name: `Karen`,
        surname: `Jones`,
        age: 92,
        email: `tangerinetornado@email.com`,
        premiumAccount: true,
        role: `user`,
    },
    {
        name: `Barbara`,
        surname: `Miller`,
        age: 77,
        email: `dreamydolphin@email.com`,
        premiumAccount: false,
        role: `admin`,
    },
];