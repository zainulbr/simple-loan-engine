CREATE OR REPLACE FUNCTION check_total_investment()
RETURNS TRIGGER AS $$
DECLARE
    total_investment INT;
    loan_amount INT;
BEGIN
    -- Ambil total investasi yang sudah ada + investasi baru
    SELECT COALESCE(SUM(amount), 0) + NEW.amount INTO total_investment
    FROM loan.investments
    WHERE loan_id = NEW.loan_id;

    -- Ambil amount dari loan
    SELECT amount INTO loan_amount
    FROM loan.loans
    WHERE loan_id = NEW.loan_id;

    -- Jika total investasi melebihi loan amount, batalkan transaksi
    IF total_investment > loan_amount THEN
        RAISE EXCEPTION 'Total investment exceeds loan amount';
    END IF;

    if total_investment = loan_amount THEN
        UPDATE loan.loans
        SET state = 'invested'
        WHERE loan_id = NEW.loan_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_total_investment
BEFORE INSERT OR UPDATE ON loan.investments
FOR EACH ROW
EXECUTE FUNCTION check_total_investment();