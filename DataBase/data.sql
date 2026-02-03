-- authors
INSERT INTO authors (first_name, last_name, bio) VALUES
('Maya', 'Torres', 'Writes contemporary fiction and short stories.'),
('Daniel', 'Hughes', 'Focuses on technology and AI ethics.'),
('Leila', 'Khan', 'Author of historical mystery novels.'),
('Jon', 'Wright', 'Specializes in business and leadership.'),
('Priya', 'Desai', 'Poet and essayist.'),
('Nate', 'Collins', 'Writes science fiction and thrillers.'),
('Ava', 'Brooks', 'Children’s author and illustrator.'),
('Omar', 'Hassan', 'Researcher in data science and analytics.'),
('Grace', 'Kim', 'Young adult fiction writer.'),
('Hector', 'Vega', 'Writes travel and memoirs.');

-- books
INSERT INTO books (title, genres, published_at, price, stock, author_id) VALUES
('The Quiet Shore', 'Fiction,Drama', '2019-05-14 00:00:00', 14.99, 42, 1),
('Ethics of Machines', 'Technology,Non-Fiction', '2021-09-21 00:00:00', 29.50, 12, 2),
('Ashes of the Crown', 'Historical,Mystery', '2018-02-01 00:00:00', 18.75, 7, 3),
('Signals in the Dark', 'Sci-Fi,Thriller', '2023-11-03 00:00:00', 22.00, 19, 6),
('Leading with Clarity', 'Business,Leadership', '2020-03-10 00:00:00', 24.00, 15, 4),
('City of Paper', 'Poetry,Essay', '2017-08-19 00:00:00', 12.50, 30, 5),
('Starlight Protocol', 'Sci-Fi', '2022-06-12 00:00:00', 19.99, 9, 6),
('Tiny Atlas', 'Children,Adventure', '2016-04-22 00:00:00', 9.99, 50, 7),
('Data Stories', 'Technology,Data', '2021-01-05 00:00:00', 27.00, 14, 8),
('Winter Lines', 'YA,Fiction', '2019-12-02 00:00:00', 15.25, 18, 9),
('Sunset Roads', 'Travel,Memoir', '2018-10-11 00:00:00', 21.40, 11, 10),
('Glass Horizon', 'Sci-Fi,Drama', '2024-02-15 00:00:00', 23.60, 13, 6),
('Team Metrics', 'Business,Data', '2020-07-07 00:00:00', 26.80, 10, 8),
('Hidden Harbor', 'Mystery,Fiction', '2017-01-29 00:00:00', 16.90, 17, 3),
('Bright Kite', 'Children,Picture Book', '2015-09-09 00:00:00', 8.75, 60, 7);

-- addresses
INSERT INTO addresses (street, city, state, postal_code, country) VALUES
('1457 Maple Ave', 'Seattle', 'WA', '98109', 'USA'),
('88 Pine Street', 'Boston', 'MA', '02108', 'USA'),
('2100 Market St', 'San Francisco', 'CA', '94114', 'USA'),
('19 River Lane', 'Austin', 'TX', '78701', 'USA'),
('502 Oak Blvd', 'Denver', 'CO', '80202', 'USA'),
('73 Hillcrest Rd', 'Portland', 'OR', '97205', 'USA'),
('900 Lake Dr', 'Chicago', 'IL', '60611', 'USA'),
('12 Rose Ct', 'Miami', 'FL', '33130', 'USA'),
('600 Elm St', 'Raleigh', 'NC', '27601', 'USA'),
('33 Sunset Ave', 'Phoenix', 'AZ', '85004', 'USA'),
('410 Birch Pkwy', 'Nashville', 'TN', '37203', 'USA'),
('5 Harbor Way', 'San Diego', 'CA', '92101', 'USA');

-- customers
INSERT INTO customers (name, email, address_id) VALUES
('Caroline Reed', 'caroline.reed@example.com', 1),
('Marcus Hill', 'marcus.hill@example.com', 2),
('Aisha Patel', 'aisha.patel@example.com', 3),
('Liam Chen', 'liam.chen@example.com', 4),
('Sofia Alvarez', 'sofia.alvarez@example.com', 5),
('Noah Bennett', 'noah.bennett@example.com', 6),
('Ivy Sanders', 'ivy.sanders@example.com', 7),
('Ethan Park', 'ethan.park@example.com', 8),
('Zara Coleman', 'zara.coleman@example.com', 9),
('Miguel Santos', 'miguel.santos@example.com', 10),
('Olivia Grant', 'olivia.grant@example.com', 11),
('Jackson Lee', 'jackson.lee@example.com', 12);

-- orders
INSERT INTO orders (customer_id, total_price, status) VALUES
(1, 29.98, 'PAID'),
(2, 29.50, 'SHIPPED'),
(3, 18.75, 'PENDING'),
(4, 46.00, 'DELIVERED'),
(5, 22.00, 'CANCELED'),
(6, 34.99, 'PAID'),
(7, 39.98, 'DELIVERED'),
(8, 12.50, 'PAID'),
(9, 43.80, 'SHIPPED'),
(10, 21.40, 'DELIVERED'),
(11, 24.00, 'PENDING'),
(12, 32.15, 'PAID'),
(1, 27.00, 'PAID'),
(2, 16.90, 'DELIVERED'),
(3, 9.99, 'PAID'),
(4, 52.75, 'SHIPPED'),
(5, 23.60, 'PAID'),
(6, 26.80, 'DELIVERED');

-- order_items
INSERT INTO order_items (order_id, book_id, quantity) VALUES
(1, 1, 2),    -- 2 x 14.99 = 29.98
(2, 2, 1),    -- 29.50
(3, 3, 1),    -- 18.75
(4, 5, 1),    -- 24.00
(4, 8, 1),    -- 9.99
(4, 15, 1),   -- 8.75  total 42.74? wait
(5, 4, 1),    -- 22.00
(6, 7, 1),    -- 19.99
(6, 1, 1),    -- 14.99  total 34.98
(7, 12, 1),   -- 23.60
(7, 14, 1),   -- 16.90  total 40.50
(8, 6, 1),    -- 12.50
(9, 13, 1),   -- 26.80
(9, 10, 1),   -- 15.25  total 42.05
(10, 11, 1),  -- 21.40
(11, 5, 1),   -- 24.00
(12, 1, 1),   -- 14.99
(12, 10, 1),  -- 15.25  total 30.24
(13, 9, 1),   -- 27.00
(14, 14, 1),  -- 16.90
(15, 8, 1),   -- 9.99
(16, 3, 1),   -- 18.75
(16, 2, 1),   -- 29.50  total 48.25
(17, 12, 1),  -- 23.60
(18, 13, 1);  -- 26.80
