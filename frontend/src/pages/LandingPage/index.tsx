import React from 'react';
import './index.scss';
import Grid from '@material-ui/core/Grid';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import ROUTE_NAMES from '../../constants/routes';
import { Link } from 'react-router-dom';
import Logo from '../../components/Logo';
import Typography from '@material-ui/core/Typography';
import { useTranslation } from 'react-i18next';
import TransKeys from '../../i18n/keys';
import Button from '@material-ui/core/Button';
import Anchor from '@material-ui/core/Link';
import TextField from '@material-ui/core/TextField';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            height: '100vh',
            width: '100vw',
            overflow: 'hidden',
        },

        info: {
            height: '100vh',
            padding: theme.spacing(2),
            display: 'flex',
            alignItems: 'center',
        },

        infoContainer: {
            height: '100vh',
            backgroundColor: theme.palette.primary.dark,
        },

        auth: {
            height: '100vh',
            display: 'flex',
            alignItems: 'center',
            padding: theme.spacing(3),
            backgroundColor: theme.palette.primary.contrastText,
        },

        card: {
            color: theme.palette.primary.contrastText,
        },

        logo: {
            color: theme.palette.primary.dark,
            position: 'absolute',
        },

        title: {
            fontWeight: 'bold',
        },

        form: {
            '& > *': {
                marginBottom: theme.spacing(2),
            },
        },

        privacyPolicy: {
            color: theme.palette.grey[700],
        },
    }),
);

export default function LandingPage() {
    const classes = useStyles();
    const { t } = useTranslation();

    return (
        <Grid container className={classes.container}>
            <Grid item xs={8}>
                <Box
                    width="100%"
                    className={classes.infoContainer}
                    display="flex"
                    justifyContent="flex-end"
                >
                    <Box width="100%" maxWidth={1000}>
                        <Link
                            className={classes.logo}
                            to={ROUTE_NAMES.DASHBOARD}
                        >
                            <Logo />
                        </Link>
                        <Box className={classes.info}>
                            <Box
                                className={classes.card}
                                width="70%"
                                maxWidth={375}
                                marginTop="300"
                            >
                                <Typography
                                    variant="h2"
                                    className={classes.title}
                                >
                                    {t(TransKeys.LANDING_PAGE.TITLE)}
                                </Typography>
                                <Typography variant="subtitle1" gutterBottom>
                                    {t(TransKeys.LANDING_PAGE.SUB_TITLE)}
                                </Typography>
                            </Box>
                        </Box>
                    </Box>
                </Box>
            </Grid>
            <Grid item xs={4}>
                <Box width="100%" maxWidth={400}>
                    <Box p={3} display="flex" justifyContent="flex-end">
                        <Button variant="outlined" color="primary">
                            Sign in
                        </Button>
                    </Box>
                    <Box className={classes.auth}>
                        <Box width="100%">
                            <Typography variant="h5">
                                {t(TransKeys.LANDING_PAGE.SIGN_UP)}
                            </Typography>
                            <Typography variant="body2">
                                or{' '}
                                <Anchor href="#">
                                    sign in to your account
                                </Anchor>
                            </Typography>
                            <form>
                                <Box
                                    mt={2}
                                    width="100%"
                                    className={classes.form}
                                >
                                    <TextField
                                        required
                                        fullWidth
                                        size="small"
                                        label="First Name"
                                        variant="outlined"
                                        autoComplete="given-name"
                                    />
                                    <TextField
                                        required
                                        fullWidth
                                        size="small"
                                        label="Surname"
                                        variant="outlined"
                                        autoComplete="family-name"
                                    />
                                    <TextField
                                        required
                                        fullWidth
                                        size="small"
                                        label="Email"
                                        autoComplete="email"
                                        variant="outlined"
                                    />
                                    <TextField
                                        required
                                        fullWidth
                                        size="small"
                                        label="Password"
                                        type="password"
                                        autoComplete="password"
                                        variant="outlined"
                                    />

                                    <Typography
                                        variant="body2"
                                        className={classes.privacyPolicy}
                                    >
                                        This page is protected by reCAPTCHA and
                                        is subject to the Google{' '}
                                        <Anchor
                                            target="_blank"
                                            href="https://policies.google.com/privacy"
                                        >
                                            Privacy Policy
                                        </Anchor>{' '}
                                        and{' '}
                                        <Anchor
                                            target="_blank"
                                            href="https://policies.google.com/terms"
                                        >
                                            Terms of Service
                                        </Anchor>
                                        .
                                    </Typography>

                                    <Button
                                        fullWidth
                                        color="secondary"
                                        variant="contained"
                                    >
                                        Sign Up
                                    </Button>
                                </Box>
                            </form>
                        </Box>
                    </Box>
                </Box>
            </Grid>
        </Grid>
    );
}
