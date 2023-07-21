import random

def generate_random_line():
    probability = random.random()
    if probability <= 0.05:
        return "kmst"
    elif probability <= 0.1:
        return "graph"
    elif probability <= 0.8:
        x = random.randint(1, 100)
        y = random.randint(1, 100)
        w = random.uniform(0, 1)
        return f"insert {x} {y} {w:.2f}"
    else:
        x = random.randint(1, 100)
        y = random.randint(1, 100)
        w = random.uniform(0, 1)
        return f"insert {x} {y} {w:.2f}"

def main():
    try:
        N = int(input("Enter the value of N: "))
        with open(f"input_test/{N}.requests", "w") as file:
            for _ in range(N):
                line = generate_random_line()
                file.write(line + "\n")
    except ValueError:
        print("Invalid input. Please enter an integer for N.")
    except IOError:
        print("Error writing to the file.")

if __name__ == "__main__":
    main()
